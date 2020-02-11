# RPi-AHRS-golang

Raspbery Piからソケット通信で送信されたセンサデータからMadgewickフィルタを使用して姿勢推定するserver側プログラムです。

サーバー側で演算を行うことで高速に処理を行うことができます。

推定した姿勢角はUnityにソケット通信で送り、描画することができます。

Raspberry Pi側の処理は[RPi-AHRS-Python](https://github.com/tetsuzawa/RPi-AHRS-Python)のclientのプログラムを流用します。