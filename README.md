# pneutrionoutil

My [NEUTRIONO](https://studio-neutrino.com/) utilities.

## neutrino.sh

Generate .wav from .musicxml and preserve inputs and outputs in a nice way.

### Usage

Install the script.

```sh
ln -s path/to/neutrino.sh path/to/NEUTRINO/neutrino.sh
cd path/to/NEUTRINO
```

Generate a default skeleton.json.

``` sh
./neutrino.sh -n > skeleton.json
```

After changing skeleton.json and installing .musicxml to `Score`, generate .wav.

``` sh
./neutrino.sh
```
