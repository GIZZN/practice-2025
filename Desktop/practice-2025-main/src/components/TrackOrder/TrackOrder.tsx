'use client';

import { useState } from 'react';
import styles from './TrackOrder.module.css';
import { motion, AnimatePresence } from 'framer-motion';
import { FaBox, FaTruck, FaMapMarkerAlt, FaClock, FaPhoneAlt, FaInfoCircle } from 'react-icons/fa';

interface Order {
  id: string;
  date: string;
  status: 'В ОБРАБОТКЕ' | 'ДОСТАВЛЯЕТСЯ' | 'ДОСТАВЛЕН';
  address: string;
  deliveryDate: string;
  notes?: string;
  items: {
    store: string;
    link: string;
    size: string;
    color: string;
    quantity: number;
  }[];
}

const mockOrders: { [key: string]: Order } = {
  "8281": {
    id: "8281",
    date: "02.07.2025, 18:46",
    status: "В ОБРАБОТКЕ",
    address: "Nikol",
    deliveryDate: "2025-08-06 в 22:51",
    notes: "йцу",
    items: [
      {
        store: "Wildberries",
        link: "http://localhost:3000/delivery-order",
        size: "2",
        color: "2",
        quantity: 1
      }
    ]
  }
};

const getStatusInfo = (status: Order['status']) => {
  switch (status) {
    case 'В ОБРАБОТКЕ':
      return {
        icon: <FaBox />,
        color: '#4A90E2',
        description: 'Заказ принят и обрабатывается'
      };
    case 'ДОСТАВЛЯЕТСЯ':
      return {
        icon: <FaTruck />,
        color: '#F5A623',
        description: 'Заказ передан в доставку'
      };
    case 'ДОСТАВЛЕН':
      return {
        icon: <FaMapMarkerAlt />,
        color: '#7ED321',
        description: 'Заказ успешно доставлен'
      };
  };
};

export const TrackOrder = () => {
  const [orderNumber, setOrderNumber] = useState("");
  const [order, setOrder] = useState<Order | null>(null);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);

  const handleTrack = () => {
    setIsLoading(true);
    setError("");
    
    // Имитация API запроса
    setTimeout(() => {
      const cleanOrderNumber = orderNumber.replace('#', '');
      const foundOrder = mockOrders[cleanOrderNumber];
      
      if (foundOrder) {
        setOrder(foundOrder);
        setError("");
      } else {
        setOrder(null);
        setError("Заказ не найден. Проверьте номер заказа.");
      }
      setIsLoading(false);
    }, 500);
  };

  const statusInfo = order ? getStatusInfo(order.status) : null;

  return (
    <div className={styles.container}>
      <motion.div 
        className={styles.trackingBox}
        initial={{ opacity: 0, y: 20 }}
        animate={{ opacity: 1, y: 0 }}
        transition={{ duration: 0.5 }}
      >
        <h1 className={styles.title}>Отследить заказ</h1>
        <p className={styles.subtitle}>
          Введите номер заказа для получения актуальной информации о его местоположении и статусе
        </p>

        <div className={styles.searchSection}>
          <div className={styles.inputWrapper}>
            <input
              type="text"
              value={orderNumber}
              onChange={(e) => setOrderNumber(e.target.value.trim())}
              placeholder="Введите номер заказа (например, #8281)"
              className={styles.input}
            />
            <motion.button
              className={styles.button}
              onClick={handleTrack}
              disabled={isLoading || !orderNumber}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
            >
              {isLoading ? "Поиск..." : "Отследить"}
            </motion.button>
          </div>
        </div>

        <AnimatePresence>
          {error && (
            <motion.div
              className={styles.error}
              initial={{ opacity: 0, y: -10 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0 }}
            >
              <FaInfoCircle className={styles.errorIcon} />
              {error}
            </motion.div>
          )}
        </AnimatePresence>

        <AnimatePresence>
          {order && statusInfo && (
            <motion.div
              className={styles.resultContainer}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0, y: -20 }}
            >
              <div className={styles.orderHeader}>
                <h2 className={styles.orderNumber}>Заказ #{order.id}</h2>
                <span className={styles.orderStatus} style={{ color: statusInfo.color }}>
                  {order.status}
                </span>
              </div>

              <div className={styles.orderDetails}>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Дата заказа:</span>
                  <span className={styles.detailValue}>{order.date}</span>
                </div>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Адрес доставки:</span>
                  <span className={styles.detailValue}>{order.address}</span>
                </div>
                <div className={styles.detailItem}>
                  <span className={styles.detailLabel}>Дата доставки:</span>
                  <span className={styles.detailValue}>{order.deliveryDate}</span>
                </div>
                {order.notes && (
                  <div className={styles.detailItem}>
                    <span className={styles.detailLabel}>Примечания:</span>
                    <span className={styles.detailValue}>{order.notes}</span>
                  </div>
                )}
              </div>

              <div className={styles.itemsContainer}>
                <h3 className={styles.itemsTitle}>Товары в заказе</h3>
                {order.items.map((item, index) => (
                  <div key={index} className={styles.orderItem}>
                    <div className={styles.itemHeader}>
                      <span className={styles.storeName}>{item.store}</span>
                      <a href={item.link} target="_blank" rel="noopener noreferrer" className={styles.itemLink}>
                        Открыть товар
                      </a>
                    </div>
                    <div className={styles.itemDetails}>
                      <span>Размер: {item.size}</span>
                      <span>Цвет: {item.color}</span>
                      <span>Количество: {item.quantity}</span>
                    </div>
                  </div>
                ))}
              </div>

              <div className={styles.statusInfo}>
                <div className={styles.statusIcon} style={{ background: statusInfo.color }}>
                  {statusInfo.icon}
                </div>
                <p className={styles.statusDescription}>{statusInfo.description}</p>
              </div>
            </motion.div>
          )}
        </AnimatePresence>
      </motion.div>
    </div>
  );
}; 