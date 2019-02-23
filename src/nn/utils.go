package nn

import (
	"image"
	"image/png"
	"io"
	"math"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distuv"
)

func sigmoid(r, c int, x float64) float64 {
	return 1.0 / (1 + math.Exp(-1*x))
}

func sigmoidDx(x float64) float64 {
	return x * (1 - x)
}

func sigmoidPrime(m mat.Matrix) mat.Matrix {
	rows, _ := m.Dims()
	o := make([]float64, rows)
	for i := range o {
		o[i] = 1
	}
	ones := mat.NewDense(rows, 1, o)
	return multiply(m, subtract(ones, m)) // m * (1 - m)
}

func dot(m, n mat.Matrix) mat.Matrix {
	r, _ := m.Dims()
	_, c := n.Dims()
	o := mat.NewDense(r, c, nil)
	o.Product(m, n)
	return o
}

func apply(fn func(i, j int, v float64) float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Apply(fn, m)
	return o
}

func scale(s float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Scale(s, m)
	return o
}

func multiply(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.MulElem(m, n)
	return o
}

func add(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Add(m, n)
	return o
}

func subtract(m, n mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	o := mat.NewDense(r, c, nil)
	o.Sub(m, n)
	return o
}

func addScalar(i float64, m mat.Matrix) mat.Matrix {
	r, c := m.Dims()
	a := make([]float64, r*c)
	for x := 0; x < r*c; x++ {
		a[x] = i
	}
	n := mat.NewDense(r, c, a)
	return add(m, n)
}

func randomArray(size int, v float64) (data []float64) {
	dist := distuv.Uniform{
		Min: -1 / math.Sqrt(v),
		Max: 1 / math.Sqrt(v),
	}

	data = make([]float64, size)
	for i := 0; i < size; i++ {
		data[i] = dist.Rand()
	}
	return
}

func max(v []float64) (int, float64) {
	maxI := -1
	maxV := -math.MaxFloat64

	for i := 0; i < len(v); i++ {
		if v[i] > maxV {
			maxV = v[i]
			maxI = i
		}
	}

	return maxI, maxV
}

func dataFromFile(imgFile io.Reader) (pixels []float64) {
	img, err := png.Decode(imgFile)
	if err != nil {
		panic(err)
	}

	// create a grayscale image
	bounds := img.Bounds()
	gray := image.NewGray(bounds)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			var rgba = img.At(x, y)
			gray.Set(x, y, rgba)
		}
	}
	// make a pixel array
	pixels = make([]float64, len(gray.Pix))
	for i := 0; i < len(gray.Pix); i++ {
		//pixels[i] = (float64(255-gray.Pix[i]) / 255.0 * 0.99) + 0.01
		pixels[i] = (float64(gray.Pix[i]) / 255.0 * 0.99) + 0.01
	}

	return
}
