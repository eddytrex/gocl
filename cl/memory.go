package cl

/*
#cgo CFLAGS: -I CL
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

extern void go_mem_notify(cl_mem memobj, void *user_data);
static void CL_CALLBACK c_mem_notify(cl_mem memobj, void *user_data) {
	go_mem_notify(memobj, user_data);
}

static cl_int CLSetMemObjectDestructorCallback(cl_mem memobj, void *user_data){
    return clSetMemObjectDestructorCallback(memobj, c_mem_notify, user_data);
}
*/
import "C"

import (
	"unsafe"
)

type CL_mem_notify func(memobj CL_mem, user_data unsafe.Pointer)

var mem_notify CL_mem_notify

//export go_mem_notify
func go_mem_notify(memobj C.cl_mem, user_data unsafe.Pointer) {
	mem_notify(CL_mem{memobj}, user_data)
}

func CLRetainMemObject(memobj CL_mem) CL_int {
	return CL_int(C.clRetainMemObject(memobj.cl_mem))
}

func CLReleaseMemObject(memobj CL_mem) CL_int {
	return CL_int(C.clReleaseMemObject(memobj.cl_mem))
}

func CLSetMemObjectDestructorCallback(memobj CL_mem,
	pfn_notify CL_mem_notify,
	user_data unsafe.Pointer) CL_int {

	if pfn_notify != nil {
		mem_notify = pfn_notify

		return CL_int(C.CLSetMemObjectDestructorCallback(memobj.cl_mem, user_data))

	} else {
		mem_notify = nil

		return CL_int(C.clSetMemObjectDestructorCallback(memobj.cl_mem, nil, nil))
	}
}

func CLEnqueueUnmapMemObject(command_queue CL_command_queue,
	memobj CL_mem,
	mapped_ptr unsafe.Pointer,
	num_events_in_wait_list CL_uint,
	event_wait_list []CL_event,
	event *CL_event) CL_int {

	if num_events_in_wait_list != 0 && len(event_wait_list) != int(num_events_in_wait_list) {
		return CL_INVALID_VALUE
	}

	var c_event C.cl_event
	var c_ret C.cl_int

	if num_events_in_wait_list != 0 {
		var c_event_wait_list []C.cl_event
		c_event_wait_list = make([]C.cl_event, num_events_in_wait_list)
		for i := 0; i < int(num_events_in_wait_list); i++ {
			c_event_wait_list[i] = event_wait_list[i].cl_event
		}

		c_ret = C.clEnqueueUnmapMemObject(command_queue.cl_command_queue,
			memobj.cl_mem,
			mapped_ptr,
			C.cl_uint(num_events_in_wait_list),
			&c_event_wait_list[0],
			&c_event)
	} else {
		c_ret = C.clEnqueueUnmapMemObject(command_queue.cl_command_queue,
			memobj.cl_mem,
			mapped_ptr,
			0,
			nil,
			&c_event)
	}

	if event != nil {
		event.cl_event = c_event
	}

	return CL_int(c_ret)
}

func CLGetMemObjectInfo(memobj CL_mem,
	param_name CL_mem_info,
	param_value_size CL_size_t,
	param_value *interface{},
	param_value_size_ret *CL_size_t) CL_int {

	var ret C.cl_int

	if (param_value_size == 0 || param_value == nil) && param_value_size_ret == nil {
		ret = C.clGetMemObjectInfo(memobj.cl_mem,
			C.cl_mem_info(param_name),
			0,
			nil,
			nil)
	} else {
		var size_ret C.size_t

		if param_value_size == 0 || param_value == nil {
			ret = C.clGetMemObjectInfo(memobj.cl_mem,
				C.cl_mem_info(param_name),
				0,
				nil,
				&size_ret)
		} else {
			switch param_name {
			case CL_MEM_TYPE:
				var value C.cl_mem_object_type
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = CL_mem_object_type(value)
			case CL_MEM_FLAGS:
				var value C.cl_mem_flags
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = CL_mem_flags(value)
			case CL_MEM_SIZE,
				CL_MEM_OFFSET:
				var value C.size_t
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = CL_size_t(value)
			case CL_MEM_HOST_PTR:
				var value unsafe.Pointer
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = value
			case CL_MEM_MAP_COUNT,
				CL_MEM_REFERENCE_COUNT:
				var value C.cl_uint
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = CL_uint(value)
			case CL_MEM_CONTEXT:
				var value C.cl_context
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = CL_context{value}
			case CL_MEM_ASSOCIATED_MEMOBJECT:
				var value C.cl_mem
				ret = C.clGetMemObjectInfo(memobj.cl_mem,
					C.cl_mem_info(param_name),
					C.size_t(param_value_size),
					unsafe.Pointer(&value),
					&size_ret)

				*param_value = CL_mem{value}
			default:
				return CL_INVALID_VALUE
			}
		}

		if param_value_size_ret != nil {
			*param_value_size_ret = CL_size_t(size_ret)
		}
	}

	return CL_int(ret)
}