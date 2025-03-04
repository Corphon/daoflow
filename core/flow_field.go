// core/flow_field.go

package core

import (
	"math"
	"sync"
	"time"
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

// FieldType 场类型
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

// FieldSystem 场系统
type FieldSystem struct {
	mu sync.RWMutex

	// 场组件
	fields  map[string]*Field // 场集合
	unified *UnifiedField     // 统一场

	// 系统属性
	energy    float64 // 总能量
	strength  float64 // 场强度
	coupling  float64 // 耦合强度
	resonance float64 // 共振强度

	// 配置
	config *FieldConfig

	// 状态
	state struct {
		lastUpdate time.Time
		metrics    map[string]float64
	}
}

// UnifiedField 统一场结构
type UnifiedField struct {
	mu sync.RWMutex

	// 基础场属性
	dimension int       // 场维度
	gridSize  int       // 网格大小
	fields    []*Field  // 子场集合
	boundary  []float64 // 边界条件

	// 场状态
	state struct {
		strength float64            // 总场强度
		energy   float64            // 总能量
		phase    float64            // 统一相位
		coupling float64            // 耦合强度
		metrics  map[string]float64 // 统一场指标
	}

	// 量子特性
	quantum struct {
		coherence    float64 // 相干性
		entanglement float64 // 纠缠度
		correlation  float64 // 关联度
	}

	// 场组合
	scalarField *Field // 标量场
	vectorField *Field // 向量场
	tensorField *Field // 张量场

	// 配置
	config *FieldConfig
}

// --------------------------------------------
// NewFieldSystem 创建新的场系统
func NewFieldSystem(config *FieldConfig) *FieldSystem {
	if config == nil {
		config = DefaultFieldConfig()
	}

	return &FieldSystem{
		fields:  make(map[string]*Field),
		unified: NewUnifiedField(config.Dimension),
		config:  config,
	}
}

// NewUnifiedField 创建新的统一场
func NewUnifiedField(dimension int) *UnifiedField {
	uf := &UnifiedField{
		dimension: dimension,
		gridSize:  DefaultGridSize,
		fields:    make([]*Field, 0),
		boundary:  make([]float64, dimension*2),
	}

	// 初始化子场
	uf.scalarField = NewField(ScalarField, dimension)
	uf.vectorField = NewField(VectorField, dimension)
	uf.tensorField = NewField(TensorField, dimension)

	// 添加到场集合
	uf.fields = append(uf.fields, uf.scalarField, uf.vectorField, uf.tensorField)

	// 初始化状态
	uf.state.metrics = make(map[string]float64)

	return uf
}

// GetEnergy 获取统一场总能量
func (uf *UnifiedField) GetEnergy() float64 {
	uf.mu.RLock()
	defer uf.mu.RUnlock()
	return uf.state.energy
}

// GetStrength 获取统一场强度
func (uf *UnifiedField) GetStrength() float64 {
	uf.mu.RLock()
	defer uf.mu.RUnlock()
	return uf.state.strength
}

// Evolve 演化统一场
func (uf *UnifiedField) Evolve() error {
	uf.mu.Lock()
	defer uf.mu.Unlock()

	// 演化各个子场
	for _, field := range uf.fields {
		if err := field.Evolve(); err != nil {
			return err
		}
	}

	// 更新统一场状态
	uf.updateState()

	return nil
}

// 更新统一场状态
func (uf *UnifiedField) updateState() {
	// 计算总能量
	totalEnergy := 0.0
	for _, field := range uf.fields {
		totalEnergy += field.GetStrength()
	}
	uf.state.energy = totalEnergy

	// 更新其他状态
	uf.state.strength = totalEnergy / float64(len(uf.fields))
	uf.updateQuantumProperties()
}

// 更新量子特性
func (uf *UnifiedField) updateQuantumProperties() {
	// 计算相干性
	coherence := 0.0
	for _, field := range uf.fields {
		coherence += field.GetCoherence()
	}
	uf.quantum.coherence = coherence / float64(len(uf.fields))

}

// GetEnergy 获取场系统总能量
func (fs *FieldSystem) GetEnergy() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.energy
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

	if len(f.Strength) == 0 {
		return 0
	}

	// 计算平均场强度
	total := 0.0
	count := 0
	for i := range f.Strength {
		for j := range f.Strength[i] {
			total += f.Strength[i][j]
			count++
		}
	}
	mean := total / float64(count)

	// 计算方差
	variance := 0.0
	for i := range f.Strength {
		for j := range f.Strength[i] {
			diff := f.Strength[i][j] - mean
			variance += diff * diff
		}
	}
	variance /= float64(count)

	// 根据方差计算均匀性，方差越小越均匀
	return 1.0 / (1.0 + math.Sqrt(variance))
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

// GetCoherence 获取场的相干性
func (f *Field) GetCoherence() float64 {
	f.mu.RLock()
	defer f.mu.RUnlock()

	if len(f.Strength) == 0 {
		return 0
	}

	// 计算场分布的相干性
	coherence := 0.0
	totalPoints := 0

	// 计算相邻点之间的相干性
	for i := 0; i < len(f.Strength); i++ {
		for j := 0; j < len(f.Strength[i]); j++ {
			if i < len(f.Strength)-1 {
				// 垂直方向相干性
				coherence += math.Cos(f.Strength[i+1][j] - f.Strength[i][j])
				totalPoints++
			}
			if j < len(f.Strength[i])-1 {
				// 水平方向相干性
				coherence += math.Cos(f.Strength[i][j+1] - f.Strength[i][j])
				totalPoints++
			}
		}
	}

	if totalPoints == 0 {
		return 0
	}

	// 归一化到[0,1]区间
	return (coherence/float64(totalPoints) + 1) / 2
}

// FieldState 场状态添加方法
func (fs *FieldState) GetStrength() float64 {
	// 计算平均场强度
	total := 0.0
	count := 0

	for _, row := range fs.Strength {
		for _, value := range row {
			total += value
			count++
		}
	}

	if count == 0 {
		return 0
	}
	return total / float64(count)
}

// GetDistribution 获取场分布
func (fs *FieldState) GetDistribution() []float64 {
	// 将二维场强度分布展平为一维数组
	distribution := make([]float64, 0)

	for _, row := range fs.Strength {
		distribution = append(distribution, row...)
	}

	return distribution
}

// CalculateOverlap 计算与另一个场状态的重叠度
func (fs *FieldState) CalculateOverlap(other *FieldState) float64 {
	if len(fs.Strength) != len(other.Strength) {
		return 0
	}

	var overlap float64
	var norm1, norm2 float64

	for i := range fs.Strength {
		for j := range fs.Strength[i] {
			overlap += fs.Strength[i][j] * other.Strength[i][j]
			norm1 += fs.Strength[i][j] * fs.Strength[i][j]
			norm2 += other.Strength[i][j] * other.Strength[i][j]
		}
	}

	if norm1 == 0 || norm2 == 0 {
		return 0
	}

	return overlap / math.Sqrt(norm1*norm2)
}

// GetGradient 获取场梯度
// 返回场强在空间中的变化率
func (fs *FieldState) GetGradient() []float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	gradients := make([]float64, 0)

	// 计算场强的空间导数
	for i := 0; i < len(fs.Strength)-1; i++ {
		for j := 0; j < len(fs.Strength[i])-1; j++ {
			// x方向梯度
			dx := fs.Strength[i+1][j] - fs.Strength[i][j]
			// y方向梯度
			dy := fs.Strength[i][j+1] - fs.Strength[i][j]

			// 合成梯度
			gradient := math.Sqrt(dx*dx + dy*dy)
			gradients = append(gradients, gradient)
		}
	}

	return gradients
}

// GetCoupling 获取场耦合强度
func (fs *FieldState) GetCoupling() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	// 计算平均梯度作为耦合强度指标
	gradients := fs.GetGradient()
	if len(gradients) == 0 {
		return 0
	}

	totalGradient := 0.0
	for _, g := range gradients {
		totalGradient += math.Abs(g)
	}

	// 归一化到[0,1]区间
	return totalGradient / float64(len(gradients))
}

// GetResonance 获取场共振强度
func (fs *FieldState) GetResonance() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	// 计算频率和振幅的乘积作为共振强度
	resonance := fs.Frequency * fs.Amplitude

	// 归一化到[0,1]区间
	return math.Min(1.0, resonance/10.0) // 10.0作为归一化因子
}

// GetStrength 获取场强度
func (fs *FieldSystem) GetStrength() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.strength
}

// GetCoupling 获取场耦合强度
func (fs *FieldSystem) GetCoupling() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.coupling
}

// GetResonance 获取场共振强度
func (fs *FieldSystem) GetResonance() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.resonance
}

// GetEnergy 获取场能量
func (fs *FieldState) GetEnergy() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.Energy
}

// GetEnergyFlow 获取能量流动
func (fs *FieldState) GetEnergyFlow() float64 {
	fs.mu.RLock()
	defer fs.mu.RUnlock()
	return fs.Flow
}

// GetMetrics 获取场态指标
func (fs *FieldState) GetMetrics() map[string]interface{} {
	fs.mu.RLock()
	defer fs.mu.RUnlock()

	return map[string]interface{}{
		"strength":  fs.GetStrength(),
		"energy":    fs.GetEnergy(),
		"coupling":  fs.GetCoupling(),
		"resonance": fs.GetResonance(),
		"phase":     fs.Phase,
		"frequency": fs.Frequency,
		"amplitude": fs.Amplitude,
		"flow":      fs.Flow,
	}
}

// GetState 获取场状态自身
func (f *Field) GetState() *Field {
	return f
}
