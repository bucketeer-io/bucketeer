workspace(
    name = "bucketeer",
    managed_directories = {
        "@npm": ["node_modules"],
    },
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

###############################################################################
# Go
###############################################################################
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "f2dcd210c7095febe54b804bb1cd3a58fe8435a909db2ec04e31542631cf715c",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.31.0/rules_go-v0.31.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

load("@io_bazel_rules_go//extras:embed_data_deps.bzl", "go_embed_data_dependencies")

go_embed_data_dependencies()

go_register_toolchains(version = "1.17.2")

http_archive(
    name = "bazel_gazelle",
    sha256 = "62ca106be173579c0a167deb23358fdfe71ffa1e4cfdddf5582af26520f1c66f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.23.0/bazel-gazelle-v0.23.0.tar.gz",
    ],
)

###############################################################################
# Dependencies
###############################################################################
load("//:repositories.bzl", "go_repositories")

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

###############################################################################
# NodeJS
###############################################################################
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "ee3280a7f58aa5c1caa45cb9e08cbb8f4d74300848c508374daf37314d5390d6",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/5.5.1/rules_nodejs-5.5.1.tar.gz"],
)

http_archive(
    name = "rules_nodejs",
    sha256 = "77cbc1989562c5b2268b293573deff30984ef06b129b40c36eff764af702fe2f",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/5.5.1/rules_nodejs-core-5.5.1.tar.gz"],
)

load(
    "@build_bazel_rules_nodejs//:index.bzl",
    "node_repositories",
    "yarn_install",
)

node_repositories(
    node_version = "16.15.1",
    yarn_version = "1.22.19",
)

yarn_install(
    name = "npm-v2",
    package_json = "//ui/web-v2:package.json",
    yarn_lock = "//ui/web-v2:yarn.lock",
)

###############################################################################
# Protobuf
###############################################################################
http_archive(
    name = "com_google_protobuf",
    sha256 = "b8ab9bbdf0c6968cf20060794bc61e231fae82aaf69d6e3577c154181991f576",
    strip_prefix = "protobuf-3.18.1",
    urls = ["https://github.com/protocolbuffers/protobuf/releases/download/v3.18.1/protobuf-all-3.18.1.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

# Googleapis
http_archive(
    name = "com_github_googleapis_googleapis",
    build_file = "@//:BUILD.googleapis",
    sha256 = "963c1126e22cfbc68a96d32f40fb76fad5ba51755a61e98c8cdfdfccd6d3354a",
    strip_prefix = "googleapis-83e756a66b80b072bd234abcfe89edf459090974/",
    urls = ["https://github.com/googleapis/googleapis/archive/83e756a66b80b072bd234abcfe89edf459090974.zip"],
)

###############################################################################
# Docker
###############################################################################
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "92779d3445e7bdc79b961030b996cb0c91820ade7ffa7edca69273f404b085d5",
    strip_prefix = "rules_docker-0.20.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.20.0/rules_docker-v0.20.0.tar.gz"],
)

load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load("@io_bazel_rules_docker//repositories:deps.bzl", container_deps = "deps")

container_deps()

load("@io_bazel_rules_docker//container:container.bzl", "container_pull")

container_pull(
    name = "bucketeer-web-nginx",
    registry = "ghcr.io",
    repository = "bucketeer-io/bucketeer-web-nginx",
    tag = "1.13.11-alpine",
)

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()
