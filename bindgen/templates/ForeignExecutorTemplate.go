const uniffiRustTaskCallbackSuccess byte = 0
const uniffiRustTaskCallbackCancelled byte = 1
const uniffiForeignExecutorCallbackSuccess byte = 0
const uniffiForeignExecutorCallbackCanceled byte = 1
const uniffiForeignExecutorCallbackError byte = 2

// Encapsulates an executor that can run Rust tasks
type UniFfiForeignExecutor struct {}

var FfiConverterForeignExecutorINSTANCE = UniFfiForeignExecutor{}

func uniffiInitForeignExecutor() {
	// TODO?
}

func (c UniFfiForeignExecutor) Lower(value UniFfiForeignExecutor) C.int {
	// TODO: pointer cast & pin memory
	return 0
}

func (c UniFfiForeignExecutor) Write(writer io.Writer, value UniFfiForeignExecutor) {
	// TODO: unclear how to do this
}

func (c UniFfiForeignExecutor) Lift(value C.int) UniFfiForeignExecutor {
	// TODO: cast pointer back
	return UniFfiForeignExecutor {}
}

func (c UniFfiForeignExecutor) Read(reader io.Reader) UniFfiForeignExecutor {
	// TODO: unclear how to to this
	return UniFfiForeignExecutor {}
}


type {{ type_|ffi_destroyer_name }} struct {}

func ({{ type_|ffi_destroyer_name }}) Destroy(_ {{ type_name }}) {}

