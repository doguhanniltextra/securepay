package com.securepay.notification.consumer;

import com.securepay.notification.model.PaymentEvent;
import lombok.extern.slf4j.Slf4j;
import org.springframework.kafka.annotation.KafkaListener;
import org.springframework.stereotype.Service;

@Slf4j
@Service
public class PaymentEventConsumer {

    @KafkaListener(
        topics = "payment.initiated",
        groupId = "notification-service-group"
    )
    public void handlePaymentEvent(PaymentEvent event) {
        log.info("Payment notification received: " +
            "payment_id={}, from={}, to={}, amount={} {}",
            event.getPaymentId(),
            event.getFromAccount(),
            event.getToAccount(),
            event.getAmount(),
            event.getCurrency()
        );
        
        // Simulate operation finished (push notification, SMS, etc.)
        log.debug("Notification sent for payment_id={}", event.getPaymentId());
    }
}
