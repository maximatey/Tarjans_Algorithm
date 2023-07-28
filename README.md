# Tarjans_Algorithm

Algoritma Tarjan merupakan salah satu metode untuk menentukan Strongly Connected Components (SCCs) pada suatu graph.

## Spesifikasi / Requirements

- Go 1.20
- Python/Npm/Alternatif untuk mengaktifkan localhost
- Javascript/Typescript
- WASM

## Cara Menggunakan

Download/clone repositori ini

- Command Line Interface

  Pada `/Tarjans_Algorithm`, jalankan perintah `go run CLI/src/main.go`
  Pada file `main.go` tersebut dapat diubah jenis inputnya di fungsi utama main

```
// Uncomment the following line to read input from a file
inputFile, _ := os.Open("test/input.txt")
scanner := bufio.NewScanner(inputFile)

// Comment the following lines if reading from a file
// fmt.Println("Enter the graph edges (type 'done' to finish input):")
// scanner := bufio.NewScanner(os.Stdin)
```

Untuk input satu-per-satu, input dapat diakhiri dengan menginput "done"

- GUI

Recompile wasm file :
Pada `/Tarjans_Algorithm/cmd/wasm`, jalankan perintah

```
    GOOS=js GOARCH=wasm go build -o ../../assets/maingo.wasm
```

Menjalankan program :
Pada `/Tarjans_Algorithm`, jalankan localhost dengan python atau npm pada server 8080

```
    python -m http.server 8080
```

Setelah itu, program dapat diakses pada `http://localhost:8080/assets/`

## Penjelasan Singkat

Algoritma Tarjan adalah algoritma yang digunakan untuk mencari komponen-komponen terhubung dalam sebuah graf (graph) berarah atau tidak berarah. Algoritma ini dikembangkan oleh Robert Tarjan pada tahun 1972.

- Kompleksitas Algoritma Tarjan

Kompleksitas algoritma Tarjan untuk graf berukuran N dengan M edge adalah O(N + M). Algoritma ini menggunakan pendekatan Depth-First Search (DFS) untuk mengelola informasi tentang simpul-simpul dalam graf, mencari komponen terhubung dan mencatat titik sengketa.

- Modifikasi mendeteksi Strong Bridges

Proses pencarian jembatan (bridges) dilakukan dengan menggunakan pendekatan algoritma Tarjan yang dimodifikasi untuk mencari strong edges. Bagian utama dari algoritma ini adalah DFS (Depth-First Search) yang digunakan untuk mengunjungi setiap simpul dalam graf dan menghitung nilai atribut "LowLink" dan "Index" pada setiap simpul.

Jika ada sisi yang menghubungkan simpul dengan ancestor-nya dalam DFS tree dan memiliki nilai "LowLink" dari simpul tujuan lebih besar dari nilai "Index" dari simpul asal, maka sisi tersebut adalah sebuah jembatan (bridge).

Sisi-sisi yang terdeteksi sebagai jembatan akan dicatat dalam atribut "Bridges" pada struktur data "Graph" dan kemudian akan ditampilkan di akhir program.

- Jenis Edges

-- Back Edge, adalah sebuah edge berarah dalam graf berarah yang menghubungkan simpul dengan salah satu ancestor-nya dalam DFS tree. Back edge digunakan dalam algoritma DFS untuk mendeteksi siklus dalam graf.

-- Forward Edge, adalah sebuah edge berarah dalam graf berarah yang menghubungkan simpul dengan salah satu descendant-nya dalam DFS tree. Forward edge digunakan untuk memberikan informasi tentang hubungan lebih lanjut antara simpul-simpul dalam graf.

-- Cross Edge, adalah sebuah edge berarah dalam graf berarah yang menghubungkan simpul yang bukan ancestor atau descendant dari simpul yang sedang diproses (current node) dalam DFS tree. Cross edge memberikan informasi tentang keterkaitan antara simpul-simpul yang berbeda dalam graf.

## Framework Frontend

WASM, HTML, Javascript

## Referensi

- [_Tarjan's Algorithm to find Strongly Connected Components_](https://www.geeksforgeeks.org/tarjan-algorithm-find-strongly-connected-components/)

- [_Strongly Connected Components_](https://www.geeksforgeeks.org/strongly-connected-components/)

- [_Bridges in a Graph_](https://www.geeksforgeeks.org/bridge-in-a-graph/)

- [Tarjan's Strongly Connected Component (SCC) Algorithm (UPDATED) | Graph Theory](https://www.youtube.com/watch?v=wUgWX0nc4NY&t=517s)
