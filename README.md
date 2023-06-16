# kdb2-crawler
[W.I.P] a KdB (https://kdb.tsukuba.ac.jp) crawler for KdB2 project

それぞれのディレクトリには以下のものが格納されています:

- kdb2csv: KdBから全科目の情報が含まれるcsvを入手し、csvの内容を標準出力へ流す 
- kdbfetch: コマンドライン引数に科目コードを与えると、その科目のシラバスのHTML版を標準出力へ流す
- kdbmining: HTMLを標準入力へ与えると情報を構造化して標準出力へ流す
- sudachi: ElasticSearchでsudachiを使うための設定ファイル群
