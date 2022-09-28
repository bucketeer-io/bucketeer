load("@build_bazel_rules_nodejs//:providers.bzl", "DeclarationInfo")

def _ts_proto_library(ctx):
    srcs_files = [f for t in ctx.attr.srcs for f in t.files.to_list()]
    dts = None
    js = None
    for src_file in srcs_files:
        is_dts = src_file.short_path.endswith(".d.ts")
        if is_dts:
            dts = src_file
        else:
            js = src_file

    return struct(
        files = depset([dts]),
        typescript = struct(
            declarations = depset([dts]),
            es5_sources = depset([js]),
            es6_sources = depset([js]),
            transitive_declarations = depset([dts]),
            transitive_es5_sources = depset([js]),
            transitive_es6_sources = depset([js]),
            type_blacklisted_declarations = depset(),
        ),
        providers = [
            DeclarationInfo(
                declarations = depset([dts]),
                transitive_declarations = depset([dts]),
                type_blacklisted_declarations = depset([]),
            ),
        ],
    )

ts_proto_library = rule(
    implementation = _ts_proto_library,
    attrs = {
        "srcs": attr.label_list(
            allow_files = ["js", "d.ts"],
        ),
    },
)
