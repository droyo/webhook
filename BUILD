load("@io_bazel_rules_go//go:def.bzl", "go_binary", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["webhook.go"],
    importpath = "aqwari.net/cmd/webhook",
    visibility = ["//visibility:private"],
)

go_binary(
    name = "webhook",
    embed = [":go_default_library"],
    visibility = ["//visibility:public"],
)
