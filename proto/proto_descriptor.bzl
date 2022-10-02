def _proto_descriptor_impl(ctx):
    inputs = ctx.files.srcs + ctx.files.deps
    descriptors = ":".join([f.path for f in ctx.files.deps])
    args = ctx.actions.args()
    args.add(descriptors, format = "--descriptor_set_in=%s")
    args.add("--include_imports")
    args.add("--include_source_info")
    args.add(ctx.outputs.out.path, format = "--descriptor_set_out=%s")
    args.add_all([f.path for f in ctx.files.srcs])
    ctx.actions.run(
        inputs = inputs,
        outputs = [ctx.outputs.out],
        arguments = [args],
        executable = ctx.executable.compiler,
    )

proto_descriptor = rule(
    implementation = _proto_descriptor_impl,
    attrs = {
        "srcs": attr.label_list(mandatory = True, allow_files = True),
        "deps": attr.label_list(allow_files = True),
        "compiler": attr.label(
            executable = True,
            cfg = "exec",
            allow_files = True,
            default = Label("@com_google_protobuf//:protoc"),
        ),
    },
    outputs = {"out": "%{name}.pb"},
)
