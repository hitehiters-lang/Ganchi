package models

type Loader struct {
	Id                     int    `bun:"id,autoincrement,pk,notnull"`
	Name                   string `bun:"name,notnull"`
	PicturePathWindows     string
	PicturePathLinux       string
	MaxLiftWeight          int
	Length                 int
	Width                  int
	Height                 int
	AutoWeight             int
	LiftHeight             int
	EngineType             string
	Voltage                int
	FrontWheels            string
	RearWheels             string
	WheelAxis              string
	BrakeType              string
	LiftingCylinder        string
	ForkLength             int
	HydraulicLiftingEngine string
	LongmenFrameMaterial   int
	SteeringMode           string
	TurningRadius          int
	ChargingTime           string
	WorkingHours           string
	LiftingAngle           int
	BatteryType            string
	Price                  int
}

type ManualLoader struct {
	Id                 int    `bun:"id,autoincrement,pk,notnull"`
	Name               string `bun:"name,notnull"`
	PicturePathWindows string
	PicturePathLinux   string
	MaxLiftWeight      int
	DriveGear          string
	Length             int
	ForkLength         int
	ForkWidth          int
	BrakeType          string
	Control            string
	MaxSpeed           float64
	LiftingSpeed       int
	Price              int
}
