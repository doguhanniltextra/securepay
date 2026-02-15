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
