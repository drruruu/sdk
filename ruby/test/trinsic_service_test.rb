require "./test/test_helper"
require 'json'
require 'okapi'

class TrinsicServiceTest < Minitest::Test

  def get_library_path
    return File.expand_path(File.join(File.dirname(__FILE__ ), 'libs'))
  end

  def before_setup
    Okapi::set_library_path(get_library_path)
    Okapi::load_native_library
  end

  def test_that_it_has_a_version_number
    refute_nil ::Trinsic::VERSION
  end

  def vaccine_cert_unsigned_path
    return File.expand_path(File.join(File.dirname(__FILE__ ), "vaccination-certificate-unsigned.jsonld"))
  end

  def vaccine_cert_frame_path
    return File.expand_path(File.join(File.dirname(__FILE__ ), "vaccination-certificate-frame.jsonld"))
  end

  def test_trinsic_services_demo
    server_address = ENV["TRINSIC_SERVER_ADDRESS"]
    wallet_service = Trinsic::WalletService.new(server_address)

    # SETUP ACTORS
    # Create 3 different profiles for each participant in the scenario
    allison = wallet_service.create_wallet("")
    clinic = wallet_service.create_wallet("")
    airline = wallet_service.create_wallet("")

    # Store profile for later use
    # File.WriteAllBytes("allison.bin", allison.ToByteString().ToByteArray());

    # Create profile from existing data
    # var allison = WalletProfile.Parser.ParseFrom(File.ReadAllBytes("allison.bin"));

    # ISSUE CREDENTIAL
    # Sign a credential as the clinic and send it to Allison
    wallet_service.set_profile(clinic)
    text = File.open(self.vaccine_cert_unsigned_path).read
    credential_json = JSON.parse(text)

    credential = wallet_service.issue_credential(credential_json)

    puts "Credential: #{credential}"

    # STORE CREDENTIAL
    # Alice stores the credential in her cloud wallet.
    wallet_service.set_profile(allison)
    item_id = wallet_service.insert_item(credential)
    puts "item id = #{item_id}"

    # SHARE CREDENTIAL
    # Allison shares the credential with the venue.
    # The venue has communicated with Allison the details of the credential
    # that they require expressed as a JSON-LD frame.
    wallet_service.set_profile(allison)

    text2 = File.open(self.vaccine_cert_frame_path).read
    proof_request_json = JSON.parse(text2)

    credential_proof = wallet_service.create_proof(document_id=item_id, reveal_document=proof_request_json)

    puts "Proof: #{credential_proof}"

    # VERIFY CREDENTIAL
    # The airline verifies the credential
    wallet_service.set_profile(airline)
    valid = wallet_service.verify_proof(credential_proof)

    puts "Verification result: #{valid}"

    assert(valid, "Credential is valid!")
  end
end