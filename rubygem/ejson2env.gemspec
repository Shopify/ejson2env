# coding: utf-8
require File.expand_path('../lib/ejson2env/version', __FILE__)

files = File.read("MANIFEST").lines.map(&:chomp)

Gem::Specification.new do |spec|
  spec.name          = "ejson2env"
  spec.version       = EJSON2ENV::VERSION
  spec.authors       = ["Catherine Jones"]
  spec.email         = ["catherine.jones@shopify.com"]
  spec.summary       = %q{Decrypt EJSON secrets and export them as environment variables}
  spec.description   = %q{Decrypt EJSON secrets and export them as environment variables.}
  spec.homepage      = "https://github.com/Shopify/ejson2env"
  spec.license       = "MIT"

  spec.files         = files
  spec.executables   = ["ejson2env"]
  spec.test_files    = []
  spec.require_paths = ["lib"]
end
