package handler_test

import "time"

// Helper functions for pointers
func stringPtr(s string) *string     { return &s }
func boolPtr(b bool) *bool           { return &b }
func floatPtr(f float64) *float64    { return &f }
func timePtr(t time.Time) *time.Time { return &t }
