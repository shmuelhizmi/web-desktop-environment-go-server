package types

type SystemInformation struct {
	CPU   CPUInformation    `json:"cpu"`
	Ram   RamInformation    `json:"ram"`
	Disks []DiskInformation `json:"disks"`
	OS    OSInformation     `json:"os"`
}

type CPUInformation struct {
	Manufacturer  string `json:"manufacturer"`
	BrandName     string `json:"brandName"`
	Speed         string `json:"speed"`
	SpeedMin      string `json:"speedMin"`
	SpeedMax      string `json:"speedMax"`
	Cores         string `json:"cores"`
	PhysicalCores string `json:"physicalCores"`
}

type RamInformation struct {
	Total string `json:"total"`
	Free  string `json:"free"`
}

type DiskInformation struct {
	Name   string `json:"name"`
	Vendor string `json:"vendor"`
	Total  string `json:"total"`
}

type OSInformation struct {
	Platform string `json:"platform"`
	Kernel   string `json:"kernel"`
	Hostname string `json:"hostname"`
}
