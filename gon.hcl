source    = ["./dist/yktr-macos_darwin_amd64/yktr"]
bundle_id = "com.github.winebarrel.yktr"

apple_id {
  username = "sugawara@winebarrel.jp"
  password = "@keychain:altool"
}

sign {
  application_identity = "Developer ID Application: Genki Sugawara"
}

zip {
  output_path = "./dist/yktr_darwin_arm64.zip"
}
