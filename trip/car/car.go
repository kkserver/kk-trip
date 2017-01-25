package car

type Car struct {
	Id int64 `json:"id"`

	LicenceCode string `json:"licenceCode"` //行驶证

	Brand    string `json:"brand"`    //品牌名
	Name     string `json:"name"`     //车辆明
	PlateNo  string `json:"plateNo"`  //车牌号
	Capacity string `json:"capacity"` //排量

	Latitude  float64 `json:"latitude"`  //经度
	Longitude float64 `json:"longitude"` //纬度
	Ip        string  `json:"string"`    //IP地址

	Atime int64 `json:"atime"`
	Ctime int64 `json:"ctime"`
}
