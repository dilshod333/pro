package models

type Device struct {
	UserID              string `bson:"user_id" json:"user_id"`
	Name                string `bson:"name" json:"name"`
	DeviceType          string `bson:"device_type" json:"device_type"`
	DeviceName          string `bson:"device_name" json:"device_name"`
	DeviceStatus        string `bson:"device_status" json:"device_status"`
	ConfigurationSettings string `bson:"configuration_settings" json:"configuration_settings"`
	LastUpdated         string `bson:"last_updated" json:"last_updated"`
	Location  			string `bson:"location" json:"location"`
}


type DeleteDeviceID struct {
	DeviceId string `bson:"device_id" json:"device_id"`
}

