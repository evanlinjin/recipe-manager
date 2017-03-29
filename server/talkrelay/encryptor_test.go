package talkrelay

import "testing"

func TestEncryptor_Encrypt(t *testing.T) {
	enc := MakeEncryptor()

	key, e := enc.MakeRandomKey()
	if e != nil {
		t.Error(e)
	}
	t.Log("KEY:", string(key))

	var ptList [10][]byte
	ptList[0] = []byte("Hello, this is a test.")
	ptList[1] = []byte("This is another test!")
	ptList[2] = []byte("Yet another~")
	ptList[3] = []byte("How many?")
	ptList[4] = []byte("One with spaces at end   ")
	ptList[5] = []byte("   ")
	ptList[6] = []byte(" Okay ")
	ptList[7] = []byte("Hello, Still testing here! Do you have this?")
	ptList[8] = []byte("Hi")
	ptList[9] = []byte("~!@#$%^&*()_+")

	for i, pt := range ptList {
		ct, err := enc.Encrypt(pt)
		if err != nil {
			t.Error(err)
		}
		t.Log("[", i, "] ", string(ct))
	}
}

func TestEncryptor_Decrypt(t *testing.T) {
	enc := MakeEncryptor()

	if e := enc.ReadEncodedKey([]byte("ggw-eExRUaA37UeKIQpP8Q")); e != nil {
		t.Error(e)
	}

	var ctList [10][]byte
	ctList[0] = []byte("71_9jbofU4lmIlL5KrOfoJhzhVlev-hx48OhM8i63Fq0YDsP7hv81iN__Zof-wG5")
	ctList[1] = []byte("ZKIC2Cq1B2CeBxekhQFz0O8c0rFB8PHpG_KkKP826DVHlfngj6fLNM_Kc7TrgIaJ")
	ctList[2] = []byte("obfXZYCCYKoERqUw_gh8Ikex03Y69ErX3pL-WXpI0hw")
	ctList[3] = []byte("4zJvpseptVofSwkDBLdv666tUXQ5lFptFReafEFsuQI")
	ctList[4] = []byte("I3ikBtHsCq-_-DIrDKzPRUBZMYJkpwszoNoLHfaHb1aGKwFk2Crrkk1SFuZn-UYf")
	ctList[5] = []byte("X7zaqpT8IFshJtaTM4gHJNr77dRpq02v0IIv9qJzrfQ")
	ctList[6] = []byte("rPJxtWQ5c_WiGyv_FXm7Ke0a6kp_rrzM0hBWUgOiVMw")
	ctList[7] = []byte("oCH5BWql5OVMfH69D8nGiZ_RX9I567xADZ5Mx8G1cRXt90lDVWuPl1XIArx7ZGIz6tCRgjNMGKDRCnKaZJ9k0Q")
	ctList[8] = []byte("rDtntUtWO0fl5P_fY6vra1_HaFuY_QUa-sWxbxCyHx8")
	ctList[9] = []byte("vwZfZ6TU0Cf8_COPctntGtBPHLs0hNE4_LN8bK9tylk")

	for i, ct := range ctList {
		pt, err := enc.Decrypt(ct)
		if err != nil {
			t.Error(err)
		}
		t.Log("[", i, "] ", string(pt))
	}
}

func TestEncryptor_Decrypt2(t *testing.T) {
	enc := MakeEncryptor()

	//if e := enc.ReadEncodedKey([]byte("AyWFv1eGEtrd_JKT-emcBA")); e != nil {
	//	t.Error(e)
	//}

	pt, e := enc.Decrypt([]byte("qIqQjYAmH6dE1xGxgUbgv5PRLqPydKWrMGQvXHvPkS2-GsrcVS_kQ00ZipZUq3cd8Hm048a6wOe9YxoXyEAHZGh-mpNq6FvHrXNJw0h05yA"))
	if e != nil {
		t.Error(e)
	}
	t.Log(string(pt))
}