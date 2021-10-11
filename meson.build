project('gxs')
golang = find_program('go')
main = join_paths(meson.current_source_dir(), 'cmd/main.go')
git_ver = run_command('git', 'log', '-n', '1', '--format=%h').stdout().strip()
vers = '-X main.version=' + git_ver
flags = ['-ldflags'] + [vers] + ['-trimpath', '-buildmode=pie', '-mod=readonly', '-modcacherw']

in_files = run_command('find', meson.current_source_dir(), '-type', 'f', '-name', '*.go').stdout().strip().split()
in_files += 'go.mod'
in_files += 'go.sum'

gxs = custom_target(
    'gxs',
    output: 'gxs',
    input: in_files,
    build_by_default: true,
    command: [ golang, 'build', flags, '-o','@OUTPUT@', main],
)

test('gxs checks',
     golang,
     workdir: meson.current_source_dir(),
     args: ['test', '-v', './...'])

test('gxs files',
     find_program('make'),
     workdir: join_paths(meson.current_source_dir(), 'tests'),
     args: ['GXS=' + join_paths(meson.current_build_dir(), 'gxs')])
