package drive

/**
 * 司机
 */
type Driver struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`        //姓名
	Code        string `json:"code"`        //身份证
	LicenceCode string `json:"licenceCode"` //驾驶证
	Ctime       int64  `json:"ctime"`
}
