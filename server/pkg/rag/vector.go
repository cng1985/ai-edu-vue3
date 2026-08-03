package rag

import (
	"encoding/binary"
	"math"
)

const LocalEmbeddingDim = 256

// CosineSimilarity 计算两个向量的余弦相似度
func CosineSimilarity(a, b []float32) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// EncodeFloat32Slice 将 float32 切片序列化为字节
func EncodeFloat32Slice(v []float32) []byte {
	buf := make([]byte, 4+len(v)*4)
	binary.LittleEndian.PutUint32(buf[:4], uint32(len(v)))
	for i, f := range v {
		binary.LittleEndian.PutUint32(buf[4+i*4:], math.Float32bits(f))
	}
	return buf
}

// DecodeFloat32Slice 从字节反序列化 float32 切片
func DecodeFloat32Slice(data []byte) []float32 {
	if len(data) < 4 {
		return nil
	}
	n := int(binary.LittleEndian.Uint32(data[:4]))
	if len(data) < 4+n*4 {
		return nil
	}
	out := make([]float32, n)
	for i := 0; i < n; i++ {
		out[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[4+i*4:]))
	}
	return out
}

// LocalEmbed 本地哈希嵌入（无需 API，适合离线/开发环境）
func LocalEmbed(text string) []float32 {
	vec := make([]float32, LocalEmbeddingDim)
	runes := []rune(text)
	for i := 0; i < len(runes); i++ {
		for l := 1; l <= 3 && i+l <= len(runes); l++ {
			gram := string(runes[i : i+l])
			h := fnvHash(gram)
			idx := h % LocalEmbeddingDim
			sign := float32(1)
			if h&1 == 1 {
				sign = -1
			}
			vec[idx] += sign
		}
	}
	normalize(vec)
	return vec
}

func fnvHash(s string) uint32 {
	h := uint32(2166136261)
	for i := 0; i < len(s); i++ {
		h ^= uint32(s[i])
		h *= 16777619
	}
	return h
}

func normalize(v []float32) {
	var sum float64
	for _, f := range v {
		sum += float64(f) * float64(f)
	}
	if sum == 0 {
		return
	}
	norm := float32(math.Sqrt(sum))
	for i := range v {
		v[i] /= norm
	}
}
