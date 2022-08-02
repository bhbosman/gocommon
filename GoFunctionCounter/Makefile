all: DeleteAllGeneratedFiles IDataShutDown

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go



IDataShutDown:
	mockgen -package GoFunctionCounter -generateWhat ddd -destination IGoFunctionCounter_InterfaceMethods.go . IGoFunctionCounter

