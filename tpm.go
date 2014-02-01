// tpm provides basic access to a hardware TPM module.
package tpm

// #cgo LDFLAGS: -ltspi
// #include <trousers/tss.h>
// #include <trousers/trousers.h>
// #include <stdlib.h>
// #include <string.h>
import "C"

import (
	"errors"
	"unsafe"
)

// TPMContext contains the internal structures for using the TPM.
type TPMContext struct {
	ctx C.TSS_HCONTEXT
	tpm C.TSS_HTPM
	srk C.TSS_HKEY
}

// These errors are defined with the error messages from the TSS library.
var (
	ErrBadParameter  = errors.New("tpm: bad parameter")
	ErrInvalidHandle = errors.New("tpm: invalid handle")
	ErrInternalError = errors.New("tpm: internal error")
	ErrUnknown       = errors.New("tpm: unknown error")
)

// NewTPMContext initialises a new TPM context and sets up the
// internal TSS structures.
func NewTPMContext() (*TPMContext, error) {
	var result C.TSS_RESULT
	var ctx TPMContext

	C.Tspi_Context_Create(&ctx.ctx)
	result = C.Tspi_Context_Connect(ctx.ctx, nil)
	if result != C.TSS_SUCCESS {
		switch result {
		case C.TSS_E_INVALID_HANDLE:
			return nil, ErrInvalidHandle
		case C.TSS_E_INTERNAL_ERROR:
			return nil, ErrInternalError
		default:
			return nil, ErrUnknown
		}
	}
	result = C.Tspi_Context_GetTpmObject(ctx.ctx, &ctx.tpm)
	if result != C.TSS_SUCCESS {
		switch result {
		case C.TSS_E_INVALID_HANDLE:
			return nil, ErrInvalidHandle
		case C.TSS_E_INTERNAL_ERROR:
			return nil, ErrInternalError
		case C.TSS_E_BAD_PARAMETER:
			return nil, ErrBadParameter
		default:
			return nil, ErrUnknown

		}
	}
	return &ctx, nil
}

// Destroy properly shuts down the TPM context.
func (ctx *TPMContext) Destroy() error {
	result := C.Tspi_Context_FreeMemory(ctx.ctx, nil)
	if result != C.TSS_SUCCESS {
		switch result {
		case C.TSS_E_INVALID_HANDLE:
			return ErrInvalidHandle
		case C.TSS_E_INTERNAL_ERROR:
			return ErrInternalError
		default:
			return ErrUnknown
		}
	}
	return nil
}

func (ctx *TPMContext) Random(n uint32) (rdata []byte, err error) {
	var randbytes *C.BYTE
	result := C.Tspi_TPM_GetRandom(ctx.tpm, C.UINT32(n), &randbytes)
	if result != C.TSS_SUCCESS {
		switch result {
		case C.TSS_E_INVALID_HANDLE:
			return nil, ErrInvalidHandle
		case C.TSS_E_INTERNAL_ERROR:
			return nil, ErrInternalError
		case C.TSS_E_BAD_PARAMETER:
			return nil, ErrBadParameter
		default:
			return nil, ErrUnknown

		}
	}

	rand := C.malloc(C.size_t(n))
	if rand != nil {
		C.memcpy(rand, unsafe.Pointer(randbytes), C.size_t(n))
		rdata = C.GoBytes(rand, C.int(n))
		C.free(rand)
	}
	C.Tspi_Context_FreeMemory(ctx.ctx, randbytes)
	return

}
