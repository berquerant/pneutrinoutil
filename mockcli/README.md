# MockCLI

```
mock pneutrinoutil CLI for testing

Usage:
  mockcli [flags]

Flags:
      --debug                        enable debug
      --desc string                  description of config
      --dry                          dryrun
      --duration duration            Process duration
      --enhanceBreathiness float32   [0, 100]%
      --env strings                  names of additional environment variables to allow reading; all allows everythings
  -e, --exclude strings              exclude task names
      --fail                         If true, exit with 1
      --formantShift float32         change voice quality (default 1)
  -h, --help                         help for mockcli
      --hook string                  command to be executed after running, result dir will be passed to 1st argument
  -i, --include strings              include task names
      --inference int                quality, processing speed: 2 (elements), 3 (standard) or 4 (advanced) (default 3)
      --list-tasks                   list task names
      --model string                 singer (default "MERROW")
  -n, --neutrinoDir string           NEUTRINO directory (default "./dist/NEUTRINO")
      --parallel int                 number of parallel (default 1)
      --pitchShiftNsf float32        change pitch via NSF
      --pitchShiftWorld float32      change pitch via WORLD
      --play string                  play command generated wav after running, wav file will be passed to 1st argument
      --randomSeed int               random seed (default 1234)
      --score string                 score file, required
  -s, --shell string                 shell command to execute (default "bash")
      --smoothFormant float32        [0, 100]%
      --smoothPitch float32          [0, 100]%
      --styleShift int               change the key and estimate to change the style of singing
      --thread int                   number of parallel in session (default 4)
  -w, --workDir string               working directory; $HOME/.pneutrinoutil or .pneutrinoutil if no $HOME (default "/Users/sin/.pneutrinoutil")
```
