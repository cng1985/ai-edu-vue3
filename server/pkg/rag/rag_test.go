package rag

import (
	"testing"
)

func TestSplitByHeadings(t *testing.T) {
	md := "# Title\n\nSome content\n\n## Section\n\nMore content"
	sections := SplitByHeadings(md)
	if len(sections) != 2 {
		t.Fatalf("expected 2 sections, got %d", len(sections))
	}
	if sections[1].Heading != "Section" {
		t.Fatalf("unexpected heading: %s", sections[1].Heading)
	}
}

func TestLocalEmbed(t *testing.T) {
	v1 := LocalEmbed("RAG 检索增强生成")
	v2 := LocalEmbed("RAG 检索增强生成")
	if len(v1) != LocalEmbeddingDim {
		t.Fatalf("expected dim %d, got %d", LocalEmbeddingDim, len(v1))
	}
	sim := CosineSimilarity(v1, v2)
	if sim < 0.99 {
		t.Fatalf("same text should have high similarity, got %f", sim)
	}
	v3 := LocalEmbed("完全不同的内容")
	sim2 := CosineSimilarity(v1, v3)
	if sim2 >= sim {
		t.Fatalf("different text should have lower similarity: %f vs %f", sim2, sim)
	}
}

func TestEncodeDecode(t *testing.T) {
	orig := []float32{0.1, 0.2, 0.3}
	data := EncodeFloat32Slice(orig)
	decoded := DecodeFloat32Slice(data)
	if len(decoded) != len(orig) {
		t.Fatalf("length mismatch")
	}
	for i := range orig {
		if decoded[i] != orig[i] {
			t.Fatalf("value mismatch at %d", i)
		}
	}
}

func TestKeywordScore(t *testing.T) {
	score := KeywordScore("RAG 是检索增强生成技术", "概述", "什么是 RAG")
	if score <= 0 {
		t.Fatal("expected positive keyword score")
	}
}
