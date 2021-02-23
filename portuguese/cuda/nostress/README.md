![](../golangnetwork-logo.png)

(c) **Cleuton Sampaio** 2019

# CUDA sem stress

A programação **CUDA** utilizando **Go** é um pouco mais complexa do que em outras linguagens. Apesar de existirem alguns pacotes excelentes, como o [**mumax**](https://godoc.org/github.com/mumax/3/cuda), a documentação é ruim, carece de exemplos e são difíceis de utilizar. 

CUDA é para **C**, então, a melhor alternativa é usar o [**Command cgo**](https://golang.org/cmd/cgo/) e invocar uma função externa com o seu **Kernel cuda**. É o que farei nesse exemplo, onde multiplico duas matrizes usando **CUDA**.

Se você quiser saber mais sobre programação **CUDA**, leia o [**meu artigo**](https://github.com/cleuton/neuraljava/tree/master/cuda). 

## Kernel

Criei um [**Kernel simples**](./maxmul.cu) que possui a função Kernel e uma helper function para ser chamada externamente. Note que eu usei **extern C** pois é assim que o **cgo** invoca funções: 

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

A função **vecmul()** é o kernel e a função **maxmul()** é a helper. Sua função é alocar memória na **GPU**, copiar os parâmetros, invocar o kernel, e copiar o resultado. Os valores são passados por referência.

## Código go

O [**programa maxmul.go**](./maxmul.go) invoca a função **helper** e exibe o resultado: 

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

Antes da importação do pacote **C**, que permite invocar funções externas em código **C** puro (extern C), eu passo a configuração do **cgo**, indicando o protótipo da função **C**, o caminho para a **lib** e o seu nome. 

Tive que criar uma **wrapper** function no código **Go** para invocar a função externa, de modo a facilitar as coisas. Ela simplesmente passa a referência para os arrays (o endereço da primeira posição) e o tamanho da matriz (neste caso 3x3 = 9). Em **CUDA** trabalhamos com matrizes *achatadas*. 

Usei o tipo **C.float** para criar **slices** contendo minhas matrizes (transformadas em vetores). Depois, invoquei a função. Note que passei o tamanho de cada linha (ou coluna). 

## Compilando

Para compilar o código **C** use o comando: 

```
nvcc --ptxas-options=-v --compiler-options '-fPIC' -o libmaxmul.so --shared maxmul.cu
```

Você precisa ter o CUDA e o driver Nvidia instalados!

Depois, é só rodar o código **Go** com o comando: 

```
go run maxmul.go
...
[19 36 16 27 41 31 28 15 24]
```

E este é o resultado da multiplicação das matrizes!