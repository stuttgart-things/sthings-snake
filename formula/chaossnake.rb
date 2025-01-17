class Chaossnake < Formula
  desc "Play snake in your terminal and wreck havoc to your Kubernetes cluster. Lol."
  homepage "https://github.com/deggja/chaossnake"

  if OS.mac?
    url "https://github.com/deggja/chaossnake/releases/download/v0.3.0/chaossnake_0.3.0_darwin_amd64.tar.gz"
    sha256 "2c70292ff95632d54e0c6a8e14d8b5ec53d548b63eb90b23161ded1cf3f0d30d"
  elsif OS.linux?
    url "https://github.com/deggja/chaossnake/releases/download/v0.3.0/chaossnake_0.3.0_linux_amd64.tar.gz"
    sha256 "5f9f13bb9dd2fa2ff5627f3fbed87cf5945a8d6734e065da432acb3c6f625f78"
  end

  def install
    bin.install "chaossnake"
  end

  test do
    system "#{bin}/chaossnake", "version"
  end
end
