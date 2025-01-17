class Serpent < Formula
  desc "Play snake in your terminal and wreck havoc to your Kubernetes cluster. Lol."
  homepage "https://github.com/deggja/serpent"

  if OS.mac?
    url "https://github.com/deggja/chaossnake/releases/download/v0.2.0/chaossnake_0.2.0_darwin_amd64.tar.gz"
    sha256 "074e6eaeb51a871f9c8eb4e21358f855673308e1ea12c61e0f473f8e750f93cb"
  elsif OS.linux?
    url "https://github.com/deggja/chaossnake/releases/download/v0.2.0/chaossnake_0.2.0_linux_amd64.tar.gz"
    sha256 "24f718845b7fbe72ea9b5cd70f3109d389a406601af40d73646672e856ceab60"
  end

  def install
    bin.install "serpent"
  end

  test do
    system "#{bin}/serpent", "version"
  end
end
