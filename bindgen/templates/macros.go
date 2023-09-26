{#/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */#}

{% macro arg_list_decl(func) %}
	{%- for arg in func.arguments() -%}
          {%- let type_ = arg.as_type() -%}
          {%- match type_ -%}
          {%- when Type::Enum { name, module_path } -%}
              {%- let e = ci.get_enum_definition(name).expect("missing cbi") -%}
              {%- if ci.is_name_used_as_error(name) -%}
                  {{ arg.name()|var_name }} *{{ arg|type_name }}
              {%- else -%}
                  {{ arg.name()|var_name }} {{ arg|type_name }}
              {%- endif -%}
          {%- else -%}
              {{ arg.name()|var_name }} {{ arg|type_name }}
          {%- endmatch -%}
		{%- if !loop.last %}, {% endif -%}
	{%- endfor -%}
{%- endmacro %}

{% macro return_type_decl(func) %}
	{%- match func.return_type() -%}
	{%- when Some with (return_type) -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) -%}
		({{ return_type|type_name }}, error)
		{%- when None -%}
		{{ return_type|type_name }}
		{%- endmatch %}
	{%- when None -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) -%}
		error
		{%- when None -%}
		{%- endmatch %}
	{%- endmatch %}
{%- endmacro %}

{% macro return_type_decl_cb(func) %}
	{%- match func.return_type() -%}
	{%- when Some with (return_type) -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) -%}
		({{ return_type|type_name }}, *{{ throws_type|type_name }})
		{%- when None -%}
		{{ return_type|type_name }}
		{%- endmatch %}
	{%- when None -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) -%}
		*{{ throws_type|type_name }}
		{%- when None -%}
		{%- endmatch %}
	{%- endmatch %}
{%- endmacro %}

{% macro ffi_call_binding(func, prefix) %}
	{%- match func.return_type() -%}
	{%- when Some with (return_type) -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) -%}
		_uniffiRV, _uniffiErr := {% call to_ffi_call(func, prefix) %}
		if _uniffiErr != nil {
			var _uniffiDefaultValue {{ return_type|type_name }}
			return _uniffiDefaultValue, _uniffiErr
		} else {
			return {{ return_type|lift_fn }}(_uniffiRV), _uniffiErr
		}
		{%- when None -%}
		return {{ return_type|lift_fn }}({% call to_ffi_call(func, prefix) %})
		{%- endmatch -%}
	{%- when None -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) -%}
		_, _uniffiErr := {% call to_ffi_call(func, prefix) %}
		return _uniffiErr
		{%- when None -%}
		{% call to_ffi_call(func, prefix) %}
		{%- endmatch -%}
	{%- endmatch -%}
{% endmacro %}

{%- macro to_ffi_call(func, prefix) -%}
	{%- match func.throws_type() %}
	{%- when Some with (e) -%}
	rustCallWithError({{ e|ffi_converter_name }}{},
	{%- else -%}
	rustCall(
	{%- endmatch %}
	{%- match func.return_type() -%}
	{%- when Some with (return_type) -%}
	func(_uniffiStatus *C.RustCallStatus) {{ return_type|ffi_type_name }} {
		return C.{{ func.ffi_func().name() }}({% call _arg_list_ffi_call(func, prefix) -%})
	})
	{%- else -%}
	func(_uniffiStatus *C.RustCallStatus) bool {
		C.{{ func.ffi_func().name() }}({% call _arg_list_ffi_call(func, prefix) -%})
		return false
	})
	{%- endmatch %}
{%- endmacro -%}

{%- macro _arg_list_ffi_call(func, prefix) %}
	{%- if !prefix.is_empty() %}
		{{ prefix }},
	{%- endif %}
	{%- for arg in func.arguments() %}
                {{- arg|lower_fn_call }}
		{%- if !loop.last %}, {% endif %}
	{%- endfor %}
	{%- if func.arguments().len() > 0 %},{% endif %} _uniffiStatus
{%- endmacro -%}

// Arglist as used in the _UniFFILib function declations.
// Note unfiltered name but ffi_type_name filters.
-#}
{%- macro arg_list_ffi_decl(func) %}
	{%- for arg in func.arguments() %}
		{{- arg.type_().borrow()|cgo_ffi_type }} {{ arg.name() -}},
	{% endfor -%}
	RustCallStatus* out_status
{%- endmacro -%}

{%- macro async_ffi_call_binding(func, prefix) -%}
	// We create a channel, that this function blocks on, until the callback sends a result on it
	done := make(chan {{ func.result_type().borrow()|future_chan_type }})
	chanHandle := cgo.NewHandle(done)
	defer chanHandle.Delete()

	rustCall(func(_uniffiStatus *C.RustCallStatus) bool {
		C.{{ func.ffi_func().name() }}(
			{%- if !prefix.is_empty() %}
			{{ prefix }},
			{%- endif %}
			{%- for arg in func.arguments() %}
			{{- arg|lower_fn_call }},
			{%- endfor %}
			FfiConverterForeignExecutorINSTANCE.Lower(UniFfiForeignExecutor {}),
			C.UniFfiFutureCallback{{ func.result_type().future_callback_param().borrow()|cgo_ffi_callback_type }}(C.{{ func.result_type().borrow()|future_callback }}),
			unsafe.Pointer(&chanHandle),
			_uniffiStatus,
		)
		return false
	})

	// wait for things to be done
        res := <- done
	
	{%- match func.return_type() -%}
	{%- when Some with (return_type) -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) %}
	 return res.val, res.err
		{%- when None %}
	 return res.val
		{%- endmatch %}
	{%- when None -%}
		{%- match func.throws_type() -%}
		{%- when Some with (throws_type) %}
	 return res.err
		{%- when None %}
         _ = res
		{%- endmatch %}
	{%- endmatch %}
{%- endmacro -%}
