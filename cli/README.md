# CLI

```
Generate .wav from .musicxml using NEUTRINO

e.g.
pneutrinoutil --neutrinoDir /path/to/NEUTRINO --workDir /path/to/install-result --score /path/to/some.musicxml

Usage:
  pneutrinoutil [CONFIG_YML|CONFIG_JSON] [flags]
  pneutrinoutil [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  info        Print system-wide information
  skeleton    Dump default config.yml
  version     Print pneutrinoutil version

Flags:
      --debug                        enable debug
      --desc string                  description of config
      --dry                          dryrun
      --enhanceBreathiness float32   [0, 100]% (before NEUTRINO v3)
      --env strings                  names of additional environment variables to allow reading; all allows everythings
  -e, --exclude strings              exclude task names
      --formantShift float32         change voice quality (before NEUTRINO v3) (default 1)
  -h, --help                         help for pneutrinoutil
      --hook string                  command to be executed after running, result dir will be passed to 1st argument
  -i, --include strings              include task names
      --inference int                quality, processing speed: 2 (elements), 3 (standard) or 4 (advanced) (before NEUTRINO v3) (default 3)
      --list-tasks                   list task names
      --model string                 singer (default "MERROW")
  -n, --neutrinoDir string           NEUTRINO directory (default "./dist/NEUTRINO")
      --parallel int                 number of parallel (before NEUTRINO v3) (default 1)
      --pitchShiftNsf float32        change pitch via NSF (before NEUTRINO v3)
      --pitchShiftWorld float32      change pitch via WORLD (before NEUTRINO v3)
      --play string                  play command generated wav after running, wav file will be passed to 1st argument
      --randomSeed int               random seed (before NEUTRINO v3) (default 1234)
      --score string                 score file, required
  -s, --shell string                 shell command to execute (default "bash")
      --smoothFormant float32        [0, 100]% (before NEUTRINO v3)
      --smoothPitch float32          [0, 100]% (before NEUTRINO v3)
      --styleShift int               change the key and estimate to change the style of singing (before NEUTRINO v3)
      --supportModel string          support singer (NEUTRINO v3)
      --thread int                   number of parallel in session (default 4)
      --transpose int                change the key and estimate (NEUTRINO v3)
  -w, --workDir string               working directory; $HOME/.pneutrinoutil or .pneutrinoutil if no $HOME (default "/Users/sin/.pneutrinoutil")

Use "pneutrinoutil [command] --help" for more information about a command.
```
