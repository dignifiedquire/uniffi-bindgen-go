{#/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */#}

// This file was autogenerated by some hot garbage in the `uniffi` crate.
// Trust me, you don't want to mess with it!

{% import "macros.go" as go %}

#include <stdbool.h>
#include <stdint.h>

// The following structs are used to implement the lowest level
// of the FFI, and thus useful to multiple uniffied crates.
// We ensure they are declared exactly once, with a header guard, UNIFFI_SHARED_H.
#ifdef UNIFFI_SHARED_H
	// We also try to prevent mixing versions of shared uniffi header structs.
	// If you add anything to the #else block, you must increment the version suffix in UNIFFI_SHARED_HEADER_V6
	#ifndef UNIFFI_SHARED_HEADER_V6
		#error Combining helper code from multiple versions of uniffi is not supported
	#endif // ndef UNIFFI_SHARED_HEADER_V6
#else
#define UNIFFI_SHARED_H
#define UNIFFI_SHARED_HEADER_V6
// ⚠️ Attention: If you change this #else block (ending in `#endif // def UNIFFI_SHARED_H`) you *must* ⚠️
// ⚠️ increment the version suffix in all instances of UNIFFI_SHARED_HEADER_V6 in this file.           ⚠️

typedef struct RustBuffer {
	int32_t capacity;
	int32_t len;
	uint8_t *data;
} RustBuffer;

typedef int32_t (*ForeignCallback)(uint64_t, int32_t, uint8_t *, int32_t, RustBuffer *);

// Task defined in Rust that Go executes
typedef void (*RustTaskCallback)(const void *, int8_t);

// Callback to execute Rust tasks using a Go routine
//
// Args:
//   executor: ForeignExecutor lowered into a uint64_t value
//   delay: Delay in MS
//   task: RustTaskCallback to call
//   task_data: data to pass the task callback
typedef int8_t (*ForeignExecutorCallback)(uint64_t, uint32_t, RustTaskCallback, void *);

typedef struct ForeignBytes {
	int32_t len;
	const uint8_t *data;
} ForeignBytes;

// Error definitions
typedef struct RustCallStatus {
	int8_t code;
	RustBuffer errorBuf;
} RustCallStatus;

// Continuation callback for UniFFI Futures
typedef void (*RustFutureContinuation)(void * , int8_t);

// ⚠️ Attention: If you change this #else block (ending in `#endif // def UNIFFI_SHARED_H`) you *must* ⚠️
// ⚠️ increment the version suffix in all instances of UNIFFI_SHARED_HEADER_V6 in this file.           ⚠️
#endif // def UNIFFI_SHARED_H

// Needed because we can't execute the callback directly from go.
void cgo_rust_task_callback_bridge_{{ config.package_name.as_ref().unwrap() }}(RustTaskCallback, const void *, int8_t);

int8_t uniffiForeignExecutorCallback{{ config.package_name.as_ref().unwrap() }}(uint64_t, uint32_t, RustTaskCallback, void*);

void uniffiFutureContinuationCallback{{ config.package_name.as_ref().unwrap() }}(void*, int8_t);

{% for func in ci.iter_ffi_function_definitions() -%}
	{%- match func.return_type() -%}{%- when Some with (type_) %}{{ type_|cgo_ffi_type }}{% when None %}void{% endmatch %} {{ func.name() }}(
	{% call go::arg_list_ffi_decl(func) %}
);

{% endfor -%}

{%- for func in self.cgo_callback_fns() %}
int32_t {{ func }}(uint64_t, int32_t, uint8_t *, int32_t, RustBuffer *);
{%- endfor %}


