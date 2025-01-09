// core/flow_field.go

package core

import (
	"math"
	"sync"
)

// 场的基本常数
const (
	// 场强度限制
	MinFieldStrength = 0.0
	MaxFieldStrength = 100.0

	// 相互作用常数
	InteractionConstant = 8.987551787e9 // 库仑常数
	GravityConstant     = 6.67430e-11   // 引力常数

	// 场的特征尺度
	DefaultGridSize = 32    // 默认场网格大小
	MinWaveLength   = 1e-10 // 最小波长(m)
	MaxWaveLength   = 1e3   // 最大波长(m)
)

// FieldType 场的类型
type FieldType uint8

const (
	ScalarField FieldType = iota // 标量场
	VectorField                  // 向量场
	TensorField                  // 张量场
)

// Field 场的基本结构
type Field struct {
	mu sync.RWMutex

	// 场的基本属性
	Type      FieldType // 场类型
	Dimension int       // 场维度
	GridSize  int       // 网格大小
	Boundary  []float64 // 边界条件

	// 场的物理量
	Strength  [][]float64  // 场强度分布
	Potential [][]float64  // 势能分布
	Gradient  [][]Vector3D // 梯度分布

	// 场的动态特性
	WaveNumber float64 // 波数
	Frequency  float64 // 频率
	Phase      float64 // 相位

	// 相互作用特性
	Coupling    float64 // 耦合强度
	Interaction float64 // 相互作用强度

	// 阴阳属性
	YinField  *Field // 阴性场
	YangField *Field // 阳性场
}

// NewField 创建新的场
func NewField(fieldType FieldType, dimension int) *Field {
	if dimension <= 0 {
		dimension = 1
	}

	field := &Field{
		Type:      fieldType,
		Dimension: dimension,
		GridSize:  DefaultGridSize,
		Boundary:  make([]float64, dimension*2), // 每个维度两个边界
		Strength:  make([][]float64, DefaultGridSize),
		Potential: make([][]float64, DefaultGridSize),
		Gradient:  make([][]Vector3D, DefaultGridSize),

		WaveNumber:  1.0,
		Frequency:   1.0,
		Phase:       0.0,
		Coupling:    0.5,
		Interaction: 0.5,
	}

	// 初始化二维数组
	for i := 0; i < DefaultGridSize; i++ {
		field.Strength[i] = make([]float64, DefaultGridSize)
		field.Potential[i] = make([]float64, DefaultGridSize)
		field.Gradient[i] = make([]Vector3D, DefaultGridSize)
	}

	return field
}

// CalculateFieldStrength 计算场强度
func (f *Field) CalculateFieldStrength(position Vector3D) float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	var strength float64

	// 根据场类型计算场强度
	switch f.Type {
	case ScalarField:
		strength = f.calculateScalarFieldStrength(position)
	case VectorField:
		strength = f.calculateVectorFieldStrength(position)
	case TensorField:
		strength = f.calculateTensorFieldStrength(position)
	}

	return math.Max(MinFieldStrength, math.Min(strength, MaxFieldStrength))
}

// calculateScalarFieldStrength 计算标量场强度
func (f *Field) calculateScalarFieldStrength(position Vector3D) float64 {
	// 使用波动方程计算标量场
	k := f.WaveNumber
	ω := f.Frequency * 2 * math.Pi
	t := float64(0) // 可以从外部传入时间参数

	// ψ(x,t) = A * sin(kx - ωt + φ)
	amplitude := 1.0
	phase := k*position.X - ω*t + f.Phase

	return amplitude * math.Sin(phase)
}

// calculateVectorFieldStrength 计算向量场强度
func (f *Field) calculateVectorFieldStrength(position Vector3D) float64 {
	// 使用类似电磁场的方法计算
	// E = k * q / r²
	r := math.Sqrt(position.X*position.X + position.Y*position.Y + position.Z*position.Z)
	if r == 0 {
		return MaxFieldStrength
	}

	charge := 1.0 // 可以设置电荷量
	return InteractionConstant * charge / (r * r)
}

// calculateTensorFieldStrength 计算张量场强度
func (f *Field) calculateTensorFieldStrength(position Vector3D) float64 {
	// 使用引力场模型
	// F = G * m1 * m2 / r²
	r := math.Sqrt(position.X*position.X + position.Y*position.Y + position.Z*position.Z)
	if r == 0 {
		return MaxFieldStrength
	}

	mass1, mass2 := 1.0, 1.0 // 可以设置质量
	return GravityConstant * mass1 * mass2 / (r * r)
}

// CalculateFieldGradient 计算场梯度
func (f *Field) CalculateFieldGradient(position Vector3D) Vector3D {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 计算各个方向的偏导数
	h := 1e-6 // 微小位移

	// ∂F/∂x
	dx := (f.CalculateFieldStrength(Vector3D{position.X + h, position.Y, position.Z}) -
		f.CalculateFieldStrength(Vector3D{position.X - h, position.Y, position.Z})) / (2 * h)

	// ∂F/∂y
	dy := (f.CalculateFieldStrength(Vector3D{position.X, position.Y + h, position.Z}) -
		f.CalculateFieldStrength(Vector3D{position.X, position.Y - h, position.Z})) / (2 * h)

	// ∂F/∂z
	dz := (f.CalculateFieldStrength(Vector3D{position.X, position.Y, position.Z + h}) -
		f.CalculateFieldStrength(Vector3D{position.X, position.Y, position.Z - h})) / (2 * h)

	return Vector3D{dx, dy, dz}
}

// ApplyYinYangSeparation 应用阴阳分离
func (f *Field) ApplyYinYangSeparation(yinRatio float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if yinRatio < 0 || yinRatio > 1 {
		return ErrInvalidParameter
	}

	yangRatio := 1 - yinRatio

	// 创建阴阳场
	f.YinField = NewField(f.Type, f.Dimension)
	f.YangField = NewField(f.Type, f.Dimension)

	// 分配场强度
	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			strength := f.Strength[i][j]
			f.YinField.Strength[i][j] = strength * yinRatio
			f.YangField.Strength[i][j] = strength * yangRatio
		}
	}

	// 设置特性
	f.YinField.Frequency *= (1 - 0.3*yinRatio)   // 阴性频率降低
	f.YangField.Frequency *= (1 + 0.3*yangRatio) // 阳性频率升高

	f.YinField.WaveNumber *= (1 - 0.2*yinRatio)   // 阴性波数降低
	f.YangField.WaveNumber *= (1 + 0.2*yangRatio) // 阳性波数升高

	return nil
}

// CalculateInterference 计算场的干涉
func (f *Field) CalculateInterference(other *Field, position Vector3D) float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 计算两个场的叠加
	amplitude1 := f.CalculateFieldStrength(position)
	amplitude2 := other.CalculateFieldStrength(position)

	// 考虑相位差
	phaseDiff := f.Phase - other.Phase

	// 使用干涉公式: I = A1² + A2² + 2*A1*A2*cos(Δφ)
	return math.Pow(amplitude1, 2) + math.Pow(amplitude2, 2) +
		2*amplitude1*amplitude2*math.Cos(phaseDiff)
}

// Initialize 初始化场
func (f *Field) Initialize() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 重置所有场值
	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			f.Strength[i][j] = 0
			f.Potential[i][j] = 0
			f.Gradient[i][j] = Vector3D{0, 0, 0}
		}
	}

	// 验证参数
	if err := f.validateParameters(); err != nil {
		return NewCoreErrorWithCode(ErrInitialize, "failed to initialize field parameters")
	}

	// 重置动态特性
	f.WaveNumber = 1.0
	f.Frequency = 1.0
	f.Phase = 0.0
	f.Coupling = 0.5
	f.Interaction = 0.5

	// 清除阴阳分离
	f.YinField = nil
	f.YangField = nil

	return nil
}

// validateParameters 验证场参数
func (f *Field) validateParameters() error {
	if f.GridSize <= 0 {
		return NewCoreErrorWithCode(ErrInvalid, "invalid grid size")
	}
	if f.Dimension <= 0 {
		return NewCoreErrorWithCode(ErrInvalid, "invalid dimension")
	}
	if f.WaveNumber < 0 {
		return NewCoreErrorWithCode(ErrInvalid, "invalid wave number")
	}
	if f.Frequency < 0 {
		return NewCoreErrorWithCode(ErrInvalid, "invalid frequency")
	}
	return nil
}

// SetStrength 设置场强度
func (f *Field) SetStrength(strength float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	if strength < MinFieldStrength || strength > MaxFieldStrength {
		return NewCoreErrorWithCode(ErrRange, "field strength out of range")
	}

	// 设置整体场强度
	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			f.Strength[i][j] = strength
		}
	}

	return nil
}

// SetPhase 设置场相位
func (f *Field) SetPhase(phase float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 标准化相位到 [0, 2π)
	phase = math.Mod(phase, 2*math.Pi)
	if phase < 0 {
		phase += 2 * math.Pi
	}

	f.Phase = phase
	return nil
}

// Evolve 场演化
func (f *Field) Evolve() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 相位演化
	f.Phase += math.Pi / 4
	f.Phase = math.Mod(f.Phase, 2*math.Pi)

	// 验证演化参数
	if f.WaveNumber < MinWaveLength || f.Frequency <= 0 {
		return NewCoreErrorWithCode(ErrField, "invalid field evolution parameters")
	}

	// 更新其他场属性
	f.WaveNumber *= 0.99 // 波数衰减
	f.Frequency *= 0.99  // 频率衰减

	return nil
}

// GetStrength 获取平均场强度
func (f *Field) GetStrength() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	total := 0.0
	count := 0

	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			total += f.Strength[i][j]
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// GetUniformity 获取场的均匀性
func (f *Field) GetUniformity() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	// 计算场强度的标准差来衡量均匀性
	mean := f.GetStrength()
	variance := 0.0
	count := 0

	// 计算方差
	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			diff := f.Strength[i][j] - mean
			variance += diff * diff
			count++
		}
	}

	if count == 0 {
		return 1.0 // 完全均匀
	}

	// 计算标准差
	stdDev := math.Sqrt(variance / float64(count))

	// 将标准差转换为均匀性指标 [0,1]
	// 标准差越小，均匀性越高
	uniformity := 1.0 / (1.0 + stdDev)

	return uniformity
}

// Update 使用给定的相位更新场
func (f *Field) Update(phase float64) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 更新相位
	f.Phase = phase

	// 更新场强度分布
	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			// 使用波动方程更新场强度
			position := Vector3D{
				X: float64(i) / float64(f.GridSize),
				Y: float64(j) / float64(f.GridSize),
				Z: 0,
			}
			f.Strength[i][j] = f.calculateScalarFieldStrength(position)
		}
	}

	// 更新动态特性
	f.WaveNumber *= 0.99 // 波数衰减
	f.Frequency *= 0.99  // 频率衰减

	return nil
}

// Reset 重置场到初始状态
func (f *Field) Reset() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	// 重置所有场值
	for i := 0; i < f.GridSize; i++ {
		for j := 0; j < f.GridSize; j++ {
			f.Strength[i][j] = 0
			f.Potential[i][j] = 0
			f.Gradient[i][j] = Vector3D{0, 0, 0}
		}
	}

	// 重置动态特性
	f.WaveNumber = 1.0
	f.Frequency = 1.0
	f.Phase = 0.0
	f.Coupling = 0.5
	f.Interaction = 0.5

	// 清除阴阳分离
	f.YinField = nil
	f.YangField = nil

	return nil
}
