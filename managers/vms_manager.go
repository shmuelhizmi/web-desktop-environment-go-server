package managers

import (
	"errors"
	"github.com/d5/tengo"
	"github.com/d5/tengo/stdlib"
	react_fullstack_go_server "github.com/shmuelhizmi/react-fullstack-go-server"
	"github.com/shmuelhizmi/web-desktop-environment-go-server/types"
)

func GetSTDByPermissions(permissions []types.VMPermission) *tengo.ModuleMap {
	strPermissions := make([]string, 0, len(permissions))
	for _, permission := range permissions {
		strPermissions = append(strPermissions, string(permission))
	}
	return stdlib.GetModuleMap(strPermissions...)
}

func GetUInt16(obj tengo.Object) (error, uint16) {
	switch valueType := obj.(type) {
	case *tengo.Int:
		{
			return nil, uint16(valueType.Value)
		}
	}
	return errors.New("could not convert value to uint16"), 0
}
func GetString(obj tengo.Object) (error, string) {
	switch valueType := obj.(type) {
	case *tengo.String:
		{
			return nil, valueType.Value
		}
	}
	return errors.New("could not convert value to string"), ""
}

func GetViewUuid(obj tengo.Object) (error, string) {
	switch valueType := obj.(type) {
	case *tengo.Map:
		{
			id := valueType.Value["id"]
			switch idType := id.(type) {
			case *tengo.String:
				{
					return nil, idType.Value

				}
			}
		}
	}
	return errors.New("could not convert value to id string"), ""
}

func ValueifyObject(obj tengo.Object) interface{} {
	switch valueType := obj.(type) {
	case *tengo.Map:
		{
			value := valueType.Value
			resultMap := make(map[string]interface{})
			for paramName, param := range value {
				resultMap[paramName] = ValueifyObject(param)
			}
			return resultMap
		}
	case *tengo.Array:
		{
			value := valueType.Value
			resultArray := make([]interface{}, 0)
			for _, param := range value {
				resultArray = append(resultArray, ValueifyObject(param))
			}
			return resultArray
		}
	case *tengo.String:
		{
			return valueType.Value
		}
	case *tengo.Bool:
		{
			return valueType.Equals(tengo.TrueValue)
		}
	case *tengo.Int:
		{
			return valueType.Value
		}
	case *tengo.Float:
		{
			return valueType.Value
		}
	}
	return nil
}

func CreateVMSManager(dependencies types.VMSManagerDependencies) types.VMSManager {
	return types.VMSManager{
		RunVMApp: func(script string, permissions []types.VMPermission, appParams *react_fullstack_go_server.ComponentParams) error {
			tengoScript := tengo.NewScript([]byte(script))
			imports := GetSTDByPermissions(permissions)
			views := make(map[string]*react_fullstack_go_server.View) // [uuid] view
			appModule := map[string]tengo.Object{
				"view": &tengo.UserFunction{
					Name: "view",
					Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
						getLayerError, layer := GetUInt16(args[0])
						if getLayerError != nil {
							return nil, tengo.ErrInvalidArgumentType{
								Name:     "first",
								Expected: "int",
								Found:    args[0].TypeName(),
							}
						}
						getViewNameError, viewName := GetString(args[1])
						if getViewNameError != nil {
							return nil, tengo.ErrInvalidArgumentType{
								Name:     "second",
								Expected: "string",
								Found:    args[1].TypeName(),
							}
						}
						var view react_fullstack_go_server.View
						if len(args) > 2 {
							getParentViewUuidError, viewParentUuid := GetViewUuid(args[2])
							if getParentViewUuidError != nil {
								return nil, tengo.ErrInvalidArgumentType{
									Name:     "third",
									Expected: "view",
									Found:    args[2].TypeName(),
								}
							}
							viewParent, found := views[viewParentUuid]
							if found {
								view = appParams.View(layer, viewName, viewParent)
							}
						} else {
							view = appParams.View(layer, viewName, nil)
						}
						views[view.Uuid] = &view
						return &tengo.Map{
							Value: map[string]tengo.Object{
								"id": &tengo.String{
									Value: view.Uuid,
								},
								"updateParam": &tengo.UserFunction{
									Name: "updateParam",
									Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
										getNameError, paramName := GetString(args[0])
										if getNameError != nil {
											return nil, tengo.ErrInvalidArgumentType{
												Name:     "first",
												Expected: "string",
												Found:    args[1].TypeName(),
											}
										}
										paramValue := ValueifyObject(args[1])
										view.Params[paramName] = paramValue
										return tengo.TrueValue, nil
									},
								},
								"handleFunction": &tengo.UserFunction{
									Name: "handleFunction",
									Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
										getNameError, functionName := GetString(args[0])
										if getNameError != nil {
											return nil, tengo.ErrInvalidArgumentType{
												Name:     "first",
												Expected: "string",
												Found:    args[1].TypeName(),
											}
										}
										view.On(functionName, func(funcArgs ...interface{}) interface{} {
											if args[1] != nil && args[1].CanCall() {
												parsedFunctionArguments := make([]tengo.Object, 0)
												for _, argument := range funcArgs {
													newObject, _ := tengo.FromInterface(argument)
													parsedFunctionArguments = append(parsedFunctionArguments, newObject)
												}
												result, executeError := args[1].Call(parsedFunctionArguments...)
												if executeError != nil {
													return nil
												} else {
													return ValueifyObject(result)
												}
											}
											return nil
										})
										return tengo.TrueValue, nil
									},
								},
								"update": &tengo.UserFunction{
									Name: "update",
									Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
										view.Update()
										return tengo.TrueValue, nil
									},
								},
								"start": &tengo.UserFunction{
									Name: "start",
									Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
										view.Start()
										return tengo.TrueValue, nil
									},
								},
								"stop": &tengo.UserFunction{
									Name: "stop",
									Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
										view.Stop()
										return tengo.TrueValue, nil
									},
								},
							},
						}, nil
					},
				},
			}
			imports.AddBuiltinModule("app", appModule)
			tengoScript.SetImports(imports)
			tengoScript.Run()
			return nil
		},
	}
}
