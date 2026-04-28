# MGML

A high-performance Golang port of the popular responsive email framework [MJML](https://github.com/mjmlio/mjml).

The main goal of this project is to act as a drop-in replacement for the original JS library. The project strives for maximum behavioral compatibility with the original, going as far as exactly reproducing some of its architectural quirks and issues ("bug-for-bug compatibility").

**The current codebase is synchronized and corresponds to the original MJML version 4.18.0**

## Key Benefits

* **Maximum compatibility:** Almost 100% match in functionality, behavior, and final HTML output compared to the original renderer.
* **Speed and efficiency:** Significantly faster template rendering and minimal memory footprint, thanks to the compiled nature of Golang.
* **Ease of upstream synchronization:** The code isn't just rewritten; it architecturally mirrors the JS version. This makes porting new features, changes, and bug fixes from the upstream MJML repository quick and straightforward.
* **Extensibility:** The component architecture closely resembles the original. If you have custom 3rd-party MJML components written in JavaScript, rewriting them in Go will be a simple, intuitive, and predictable process.

## Performance

To evaluate the performance of mgml against other implementations, a benchmark was conducted using a modified version of the [benchmarking script](https://github.com/preslavrachev/gomjml/blob/main/bench-austin.sh) from the `gomjml` project. All tests were executed by rendering the standard [austin](https://mjml.io/try-it-live/templates/austin) template.

### Tested Versions:

* **MJML (JS):** 4.18.0 (with `--config.beautify=false`)
* **mrml-cli:** 1.7.1
* **gomjml:** v0.11.0
* **mgml:** v0.1.0

### Benchmark Results:

```
| Tool        | 100x Total (ms)     | Avg (ms)    | Max RAM (MB)| Avg CPU (%) |
|-------------|---------------------|-------------|-------------|-------------|
| gomjml      |                 427 |           4 |           1 |           0 |
| mrml        |                 160 |           1 |           2 |           0 |
| mgml        |                 469 |           4 |           1 |           0 |
| mjml (JS)   |               14912 |         149 |          83 |         8.2 |
```

## Omitted Features

To maintain full transparency, please note that the following features from the original MJML library have not been included in this port:

* **Minification and Beautification:** The `minify` and `beautify` processing steps are not implemented. These operations are not considered core functionalities of the MJML rendering framework itself. If you require minified or beautifully formatted HTML output, it is highly recommended to pass the generated HTML through dedicated external tools or libraries better suited for those specific tasks.

## Similar Projects and Alternatives

If this project does not fit your specific needs, or you are looking for implementations in other languages, you might want to check out these alternatives:

* [MJML](https://github.com/mjmlio/mjml) - The original JavaScript/Node.js framework by the mjml.io team.
* [MRML](https://github.com/jdrouet/mrml) - A popular and high-performance reimplementation written in Rust.
* [gomjml](https://github.com/preslavrachev/gomjml) - Another alternative implementation of MJML written in Golang.
* [MJML.NET](https://github.com/SebastianStehle/mjml-net) - A native MJML implementation designed for the .NET platform.

## License

This project is licensed under the [MIT License](LICENSE).

### Acknowledgements / Credits

This project is a Golang port of the original MJML framework created by [Mailjet SAS](https://mjml.io).

For detailed licensing information, full copyright notices, and the original license texts for these libraries, please refer to the [`THIRD_PARTY_NOTICES.md`](THIRD_PARTY_NOTICES.md) file.

