package responseType

type devicesResponseObject struct {
	Devices []struct {
		EasIds struct {
		} `json:"EasIds"`
		TimeZone           string `json:"TimeZone"`
		Udid               string `json:"Udid"`
		SerialNumber       string `json:"SerialNumber"`
		MacAddress         string `json:"MacAddress"`
		Imei               string `json:"Imei"`
		EasID              string `json:"EasId"`
		AssetNumber        string `json:"AssetNumber"`
		DeviceFriendlyName string `json:"DeviceFriendlyName"`
		DeviceReportedName string `json:"DeviceReportedName"`
		LocationGroupID    struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
			UUID string `json:"Uuid"`
		} `json:"LocationGroupId"`
		LocationGroupName string `json:"LocationGroupName"`
		UserID            struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
			UUID string `json:"Uuid"`
		} `json:"UserId"`
		UserName             string `json:"UserName"`
		DataProtectionStatus int    `json:"DataProtectionStatus"`
		UserEmailAddress     string `json:"UserEmailAddress"`
		Ownership            string `json:"Ownership"`
		PlatformID           struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
		} `json:"PlatformId"`
		Platform string `json:"Platform"`
		ModelID  struct {
			ID struct {
				Value int `json:"Value"`
			} `json:"Id"`
			Name string `json:"Name"`
		} `json:"ModelId"`
		Model                            string `json:"Model"`
		OperatingSystem                  string `json:"OperatingSystem"`
		PhoneNumber                      string `json:"PhoneNumber"`
		LastSeen                         string `json:"LastSeen"`
		EnrollmentStatus                 string `json:"EnrollmentStatus"`
		ComplianceStatus                 string `json:"ComplianceStatus"`
		CompromisedStatus                bool   `json:"CompromisedStatus"`
		LastEnrolledOn                   string `json:"LastEnrolledOn"`
		LastComplianceCheckOn            string `json:"LastComplianceCheckOn"`
		LastCompromisedCheckOn           string `json:"LastCompromisedCheckOn"`
		IsSupervised                     bool   `json:"IsSupervised"`
		VirtualMemory                    int    `json:"VirtualMemory"`
		OEMInfo                          string `json:"OEMInfo"`
		IsDeviceDNDEnabled               bool   `json:"IsDeviceDNDEnabled"`
		IsDeviceLocatorEnabled           bool   `json:"IsDeviceLocatorEnabled"`
		IsCloudBackupEnabled             bool   `json:"IsCloudBackupEnabled"`
		IsActivationLockEnabled          bool   `json:"IsActivationLockEnabled"`
		IsNetworkTethered                bool   `json:"IsNetworkTethered"`
		BatteryLevel                     string `json:"BatteryLevel"`
		IsRoaming                        bool   `json:"IsRoaming"`
		SystemIntegrityProtectionEnabled bool   `json:"SystemIntegrityProtectionEnabled"`
		ProcessorArchitecture            int    `json:"ProcessorArchitecture"`
		TotalPhysicalMemory              int    `json:"TotalPhysicalMemory"`
		AvailablePhysicalMemory          int    `json:"AvailablePhysicalMemory"`
		OSBuildVersion                   string `json:"OSBuildVersion"`
		DeviceCellularNetworkInfo        []struct {
			CarrierName string `json:"CarrierName"`
			CardID      string `json:"CardId"`
			PhoneNumber string `json:"PhoneNumber"`
			DeviceMCC   struct {
				Simmcc     string `json:"SIMMCC"`
				CurrentMCC string `json:"CurrentMCC"`
			} `json:"DeviceMCC"`
			IsRoaming bool `json:"IsRoaming"`
		} `json:"DeviceCellularNetworkInfo"`
		EnrollmentUserUUID string `json:"EnrollmentUserUuid"`
		ManagedBy          int    `json:"ManagedBy"`
		WifiSsid           string `json:"WifiSsid"`
		ID                 struct {
			Value int `json:"Value"`
		} `json:"Id"`
		UUID              string `json:"Uuid"`
		ComplianceSummary struct {
			DeviceCompliance []struct {
				CompliantStatus     bool   `json:"CompliantStatus"`
				PolicyName          string `json:"PolicyName"`
				PolicyDetail        string `json:"PolicyDetail"`
				LastComplianceCheck string `json:"LastComplianceCheck"`
				NextComplianceCheck string `json:"NextComplianceCheck"`
				ActionTaken         []struct {
					ActionType int `json:"ActionType"`
				} `json:"ActionTaken"`
				ID struct {
					Value int `json:"Value"`
				} `json:"Id"`
				UUID string `json:"Uuid"`
			} `json:"DeviceCompliance"`
		} `json:"ComplianceSummary,omitempty"`
	} `json:"Devices"`
	Page     int `json:"Page"`
	PageSize int `json:"PageSize"`
	Total    int `json:"Total"`
}
