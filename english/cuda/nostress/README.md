![](./golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2019

# CUDA no-stress programming

Programming **CUDA** using **Go** is a bit more complex than in other languages. Although there are some excellent packages, such as [**mumax**](https://godoc.org/github.com/mumax/3/cuda), the documentation is poor, lacks examples and it's difficult to use.

CUDA is for **C**, so the best alternative is to use [**Command cgo**](https://golang.org/cmd/cgo/) and invoke an external function with your **Cuda Kernel**. This is what I will do in this example, where I multiply two matrices using **CUDA**.

If you want to know more about **CUDA** programming, read the [**my article**](https://github.com/cleuton/neuraljava/tree/master/cuda).

## Kernel

I created a [**Simple Kernel**](./maxmul.cu) that has the Kernel function and a helper function to be called externally. Note that I used **extern C** because this is how **cgo** invokes functions:

```
#include <stdio.h>
#include <cuda.h>
 

__global__ void vecmul(float *A, float* B, float *C, int size)
{
    // Row and Column indexes: 
    int row = blockIdx.y*blockDim.y+threadIdx.y;
    int col = blockIdx.x*blockDim.x+threadIdx.x;

    // Are they bellow the maximum?
    if (col < size && row < size) {
       float result = 0;
       for(int ix=0;ix<size;ix++) {
          result += A[row*size+ix]*B[ix*size+col];
       }
       C[row*size+col] = result;
    }
}

extern "C" {

    void maxmul(float *A, float* B, float *C, int size) {

        int total = size*size;

        // Allocate device memory:
        float* gpu_A;
        float* gpu_B;
        float* gpu_C;
        int msize = total * sizeof(float);
        cudaMalloc((void**)&gpu_A, msize);
        cudaMemcpy(gpu_A,A,msize,cudaMemcpyHostToDevice);
        cudaMalloc((void**)&gpu_B, msize);
        cudaMemcpy(gpu_B,B,msize,cudaMemcpyHostToDevice);
        cudaMalloc((void**)&gpu_C,msize);

        // Blocks & grids:
        dim3 blocks(size,size);
        dim3 grid(1,1);

        // Call the kernel:
        vecmul<<<grid,blocks>>>(gpu_A,gpu_B,gpu_C,size);

        // Get the result Matrix:
        cudaMemcpy(C,gpu_C,msize,cudaMemcpyDeviceToHost);

        //Free device matrices
        cudaFree(gpu_A);
        cudaFree(gpu_B);
        cudaFree(gpu_C);
    }

}

```

The **vecmul()** function is the kernel and the **maxmul()** function is the helper. Its function is to allocate memory in the **GPU**, copy the parameters, invoke the kernel, and copy the result. Values ​​are passed by reference.

## Go code

[**Program maxmul.go**](./maxmul.go) invokes the **helper** function and displays the result: 

```
package main

/*
void maxmul(float *A, float* B, float *C, int size);
#cgo LDFLAGS: -L. -L./ -lmaxmul
*/
import "C"

import "fmt"

func Maxmul(a []C.float, b []C.float, c []C.float, size int) {
	C.maxmul(&a[0], &b[0], &c[0], C.int(size))
}

func main() {
	//in := []C.float{1.23, 4.56}
    //C.test(&in[0]) // C 1.230000 4.560000
	a := []C.float{-1,2,4,0,5,3,6,2,1}
	b := []C.float{3,0,2,3,4,5,4,7,2}
	var c []C.float = make([]C.float, 9)
	Maxmul(a,b,c,3)
	fmt.Println(c)
}
```

Before importing the **C** package, which allows to invoke external functions in pure **C** code (extern C), I pass the configuration of **cgo**, indicating the prototype of the function **C** , the path to **lib** and its name.

I had to create a **wrapper** function in the **Go** code to invoke the external function to make things easier. It simply passes the reference to the arrays (the address of the first position) and the array size (in this case 3x3 = 9). In **CUDA** we work with *flat* matrices.

I used the type **C.float** to create **slices** containing my arrays (transformed into vectors). Then I called the function. Note that I passed the size of each row (or column).

## Compiling

To compile the **C** code use the command:

```
nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmaxmul.so --shared maxmul.cu
```

You need to have CUDA and the Nvidia driver installed!

Then just run the **Go** code with the command:

```
go run maxmul.go
...
[19 36 16 27 41 31 28 15 24]
```

And this is the result of matrix multiplication!