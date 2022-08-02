all: DeleteAllGeneratedFiles IDataShutDown

DeleteAllGeneratedFiles:
	del *_InterfaceMethods.go



IDataShutDown:
	mockgen -package IDataShutDown -generateWhat ddd -destination IDataShutDown_InterfaceMethods.go . IDataShutDown

