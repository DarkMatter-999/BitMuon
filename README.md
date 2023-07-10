# BitMuon

BitMuon is a command-line application for downloading torrent files. It utilizes the custom MuonEngine for handling torrent file operations. BitMuon allows users to provide a path to a torrent file as a command-line argument and initiates the download process.

## Usage

```bash
bitmuon "path/to/torrent/file.torrent"
```

Replace `"path/to/torrent/file.torrent"` with the actual path to the torrent file you want to download.

## Installation

1. Clone the BitMuon repository:

    ```bash
    git clone https://github.com/DarkMatter-999/BitMuon.git
    ```

2. Build the BitMuon executable:

    ```bash
    cd bitmuon
    go build
    ```

3. Run BitMuon:

    ```bash
    ./bitmuon "path/to/torrent/file.torrent"
    ```

## Documentation

All the Documentation for the respective conponents can be found in [Docs folder](./docs)

## Dependencies

-   Go (version 1.16 or higher)
-   Bencode-go Module (automatically fetched by Go modules)

## License

This project is licensed under the BSD License. See the [LICENSE](./LICENSE) file for more details.

## Contact

For support, bug reports, and feature requests, please create an issue on the [BitMuon GitHub repository](https://github.com/DarkMatter-999/BitMuon/issues).

---
