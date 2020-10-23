package types

type IconType = string
type NativeIconType = string

const (
	ReactIconsIcon IconType = "icon"
	ImageIcon      IconType = "img"
)

const (
	NativeIconTypeAntDesign              NativeIconType = "AntDesign"
	NativeIconTypeEntypo                 NativeIconType = "Entypo"
	NativeIconTypeEvilIcons              NativeIconType = "EvilIcons"
	NativeIconTypeFeather                NativeIconType = "Feather"
	NativeIconTypeFontAwesome            NativeIconType = "FontAwesome"
	NativeIconTypeFontAwesome5           NativeIconType = "FontAwesome5"
	NativeIconTypeFontisto               NativeIconType = "Fontisto"
	NativeIconTypeFoundation             NativeIconType = "Foundation"
	NativeIconTypeIonicons               NativeIconType = "Ionicons"
	NativeIconTypeMaterialCommunityIcons NativeIconType = "MaterialCommunityIcons"
	NativeIconTypeMaterialIcons          NativeIconType = "MaterialIcons"
	NativeIconTypeOcticons               NativeIconType = "Octicons"
	NativeIconTypeSimpleLineIcons        NativeIconType = "SimpleLineIcons"
	NativeIconTypeZocial                 NativeIconType = "Zocial"
)

type Icon struct {
	Icon     string   `json:"icon"`
	IconType IconType `json:"type"`
}

type NativeIcon struct {
	Icon     string         `json:"icon"`
	IconType NativeIconType `json:"type"`
}
