import React, { useState, useEffect } from "react";
import Modal from "react-modal";
import Slider from "react-slick";
import { motion, AnimatePresence } from "framer-motion";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import "./Dashboard.css";
import { authService } from "./AuthService";

interface Product {
  id: string;
  name: string;
  description: string;
  price: string;
  category: string;
  images: string[];
  userName: string;
  userEmail: string;
}

const Dashboard: React.FC = () => {
  const [products, setProducts] = useState<Product[]>([]);
  const [filteredProducts, setFilteredProducts] = useState<Product[]>([]);
  const [selectedProduct, setSelectedProduct] = useState<Product | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [selectedCategory, setSelectedCategory] = useState<string>("all");
  const [priceRange, setPriceRange] = useState<number>(1000000);

  useEffect(() => {
    const fetchListings = async () => {
      try {
        const savedProducts = await authService.getListingsByOtheruser();
        if (savedProducts) {
          const updatedProducts: Product[] = savedProducts.map((prod) => ({
            id: String(prod.id),
            name: prod.productName,
            description: prod.productDescription,
            price: `${prod.price}$`,
            category: prod.category,
            images: prod.images.map(
              (imgObj: any) =>
                `data:${imgObj.contentType};base64,${imgObj.data}`
            ),
            userName: prod.userName,
            userEmail: prod.userEmail,
          }));
          setProducts(updatedProducts);
          setFilteredProducts(updatedProducts); // initially show all
        }
      } catch (error) {
        console.error("Error fetching listings:", error);
      }
    };
    fetchListings();
  }, []);

  useEffect(() => {
    const filtered = products.filter((prod) => {
      prod.price = prod.price.replaceAll("$", "");
      const price = Number(prod.price);
      const categoryMatch =
        selectedCategory === "all" || prod.category === selectedCategory;
      const priceMatch = price <= priceRange;
      prod.price = prod.price + "$";
      return categoryMatch && priceMatch;
    });
    setFilteredProducts(filtered);
  }, [products, selectedCategory, priceRange]);

  const handleProductClick = (product: Product) => {
    setSelectedProduct(product);
    setIsModalOpen(true);
  };

  const categories = Array.from(
    new Set(products.map((p) => p.category))
  ).filter(Boolean);

  const carouselSettings = {
    dots: true,
    infinite: true,
    speed: 500,
    slidesToShow: 1,
    slidesToScroll: 1,
    adaptiveHeight: true,
    arrows: true,
    responsive: [
      {
        breakpoint: 768,
        settings: {
          arrows: false,
        },
      },
    ],
  };

  return (
    <div className="app-container">
      <div className="dashboard-container">
        {/* Filter Panel */}
        <motion.div
          className="filter-panel"
          initial={{ opacity: 0, y: -10 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.4 }}
        >
          <div className="filter-group">
            <label htmlFor="category">Category</label>
            <select
              id="category"
              value={selectedCategory}
              onChange={(e) => setSelectedCategory(e.target.value)}
            >
              <option value="all">All</option>
              {categories.map((cat) => (
                <option key={cat} value={cat}>
                  {cat}
                </option>
              ))}
            </select>
          </div>

          <div className="filter-group">
            <label htmlFor="price">Max Price: ${priceRange}</label>
            <input
              type="range"
              id="price"
              min={0}
              max={1000000}
              value={priceRange}
              onChange={(e) => setPriceRange(Number(e.target.value))}
            />
          </div>
        </motion.div>

        {/* Product Grid */}
        <div className="products-grid">
          {filteredProducts.length === 0 ? (
            <div className="empty-state">No products match your filters.</div>
          ) : (
            <AnimatePresence>
              {filteredProducts.map((product) => (
                <motion.div
                  key={product.id}
                  className="product-card"
                  onClick={() => handleProductClick(product)}
                  initial={{ opacity: 0, scale: 0.95 }}
                  animate={{ opacity: 1, scale: 1 }}
                  exit={{ opacity: 0, scale: 0.95 }}
                  transition={{ duration: 0.3 }}
                >
                  <div className="product-image-container">
                    {product.images.length > 0 ? (
                      <img
                        src={product.images[0]}
                        alt={product.name}
                        className="product-thumbnail"
                      />
                    ) : (
                      <div className="no-image-placeholder">No Image</div>
                    )}
                  </div>
                  <div className="product-info">
                    <h3>{product.name}</h3>
                    <div className="price-category">
                      <span className="price">{product.price}</span>
                      <span className="category">{product.category}</span>
                    </div>
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
          )}
        </div>

        {/* Modal */}
        <Modal
          isOpen={isModalOpen}
          onRequestClose={() => setIsModalOpen(false)}
          className="product-modal"
          overlayClassName="modal-overlay"
          ariaHideApp={false}
        >
          {selectedProduct && (
            <div className="modal-content">
              <button
                className="close-button"
                onClick={() => setIsModalOpen(false)}
              >
                &times;
              </button>

              <div className="carousel-container">
                {selectedProduct.images.length > 0 ? (
                  <Slider {...carouselSettings}>
                    {selectedProduct.images.map((image, index) => (
                      <div key={index} className="slide-container">
                        <img
                          src={image}
                          alt={`${selectedProduct.name} - ${index + 1}`}
                          className="carousel-image"
                        />
                      </div>
                    ))}
                  </Slider>
                ) : (
                  <div className="no-image-placeholder">
                    No Images Available
                  </div>
                )}
              </div>

              <div className="product-details">
                <h2>{selectedProduct.name}</h2>
                <div className="meta-info">
                  <span className="price">{selectedProduct.price}</span>
                  <span className="category">{selectedProduct.category}</span>
                </div>

                <div className="scrollable-content">
                  <p className="description">{selectedProduct.description}</p>

                  <div className="contact-info">
                    <h3>Contact Seller</h3>
                    <div className="user-details">
                      <p className="seller-name">
                        <span className="icon">👤</span>
                        {selectedProduct.userName}
                      </p>
                      <p className="seller-email">
                        <span className="icon">✉️</span>
                        <a href={`mailto:${selectedProduct.userEmail}`}>
                          {selectedProduct.userEmail}
                        </a>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          )}
        </Modal>
      </div>
    </div>
  );
};

export default Dashboard;
