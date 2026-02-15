# Öğrenilenler — SecurePay

## Format
Bir hata çözüldüğünde veya önemli bir şey öğrenildiğinde buraya ekle:
- Tarih
- Ne oldu
- Nasıl çözüldü
- Bir daha nasıl önlenir

---

## [2026-02-15] Başlangıç
Henüz öğrenme kaydı yok.

## [2026-02-15] Task Durum Senkronizasyonu
Tasks.json'da "completed" görünen bir task'ın ("id: 1") dosyalarının diskte olmadığı fark edildi.
Task durumu ile dosya sistemi arasında tutarsızlık olabilir, her zaman dosya sistemini kontrol et.
Çözüm: Dosyalar yeniden oluşturuldu.

## [2026-02-15] Proto Dosyaları Eksikliği
API Gateway gRPC client implementasyonu sırasında backend servisleri için proto dosyalarının veya generated kodların olmadığı fark edildi.
Çözüm: `PaymentServiceClient` arayüzü geçici olarak mocklandı. Backend servisleri (`payment-service`, `account-service`) proto tanımları yapıldıktan sonra gerçek generated kod kullanılmalı.
Konvansiyonlarda `proto/gen/go/` klasörü beklenirken mevcut değil. But task 3.1 addressed creating proto files.

## [2026-02-15] Protoc Compiler Eksikliği ve Çözümü
Task 3.1 kapsamında proto dosyaları oluşturuldu ancak ilk başta `protoc` bulunamadığı sanıldı.
Daha sonra `go install` ile pluginler kuruldu ve `protoc` sistemde (v33.5) tespit edildi.
`proto` klasörü ayrı bir Go modülü yapıldı ve `go.work` ile `api-gateway` modülüne bağlandı.
Generated kodlar başarıyla oluşturuldu ve `grpc_clients.go` güncellendi.

## [2026-02-15] Minikube Docker Driver (WSL)
Minikube'un WSL üzerinde `docker` driver ile başlatılması gerektiği tekrarlandı. `virtualbox` veya diğer driver'lar WSL'de sorun çıkarabilir.
Çözüm: `minikube start --driver=docker` komutu kullanıldı.

## [2026-02-15] PowerShell Komut Zincirleme
PowerShell ortamında `&&` operatörü çalışmıyor.
Çözüm: Komutları ayrı ayrı çalıştırmak veya `;` (veya powershell sürümüne göre uygun operatör) kullanmak.
## [2026-02-15] Minikube ve Docker Credential Helper Sorunu (WSL)
WSL ortamında `eval $(minikube docker-env)` kullanıldığında, Docker build işlemi `docker-credential-desktop.exe: exec format error` hatası veriyor. Çünkü Minikube environment'ı Linux tabanlı olmasına rağmen Windows credential helper'ı çağırmaya çalışıyor.
Çözüm: `DOCKER_CONFIG` environment variable'ı ile credential helper içermeyen boş bir config dosyası gösterilerek build işlemi yapılabilir. Alternatif olarak imaj tag'i değiştirilip cache bypass edilebilir.

## [2026-02-15] Minikube Image Caching
Minikube, yerel olarak build edilen ve `imagePullPolicy: Never` olan imajlarda `latest` tag'ini güncellemekte zorlanıyor. Pod restart edilse bile eski imaj ID'si kullanılabiliyor.
Çözüm: Imaj tag'ini değiştirmek (örn: `v1.0.0`) en kesin çözümdür.

## [2026-02-15] SPIFFE Socket Mount Yöntemi
Kubernetes ortamında `csi.spiffe.io` driver kullanımı bazı durumlarda kararsızlık yaratabiliyor veya path sorunlarına yol açabiliyor.
Çözüm: `HostPath` volume kullanarak `/run/spire/agent-sockets` dizinini mount etmek daha stabil ve güvenilir bir yöntemdir.

## [2026-02-15] Servis Bağımlılıkları ve Başlangıç
API Gateway gibi servislerin, bağımlı oldukları backend servisleri (Payment, Account) henüz ayakta olmasa bile açılabilmesi, geliştirme ve test süreçlerini kolaylaştırır.
Çözüm: `main.go` içerisinde servis bağlantı (dial) hataları `Fatal` yerine `Warning` seviyesine çekilerek uygulamanın çökmesi engellendi ve `/health` endpointi erişilebilir kılındı.
