<template>
  <div class="login-container" :class="{ 'login-hidden': !show }">
    <!-- Background animation elements -->
    <div class="bg-shapes">
      <div class="shape shape-1"></div>
      <div class="shape shape-2"></div>
      <div class="shape shape-3"></div>
      <div class="shape shape-4"></div>
      <div class="shape shape-5"></div>
    </div>

    <!-- Main login card -->
    <div class="login-card" :class="{ 'card-hidden': !showCard }">
      <!-- Header section -->
      <div class="login-header">
        <div class="logo-container">
          <div class="logo-circle">
            <svg class="logo-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
            </svg>
          </div>
          <h2 class="login-title">Welcome Back</h2>
          <p class="login-subtitle">Sign in to your account to continue</p>
        </div>
      </div>

      <!-- Form section -->
      <form @submit.prevent="handleLogin" class="login-form">
        <!-- Email input -->
        <div class="input-group" :class="{ 'input-error': errors.email }">
          <label for="email" class="input-label">Email Address</label>
          <div class="input-container">
            <div class="input-icon">
              <svg class="input-svg" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
              </svg>
            </div>
            <input
              id="email"
              v-model="form.email"
              type="email"
              class="input-field"
              placeholder="you@example.com"
              autocomplete="email"
              required
              @focus="clearError('email')"
            />
          </div>
          <p v-if="errors.email" class="error-message">{{ errors.email }}</p>
        </div>

        <!-- Password input -->
        <div class="input-group" :class="{ 'input-error': errors.password }">
          <div class="flex justify-between items-center">
            <label for="password" class="input-label">Password</label>
            <a href="#" class="forgot-link" @click.prevent="handleForgotPassword">Forgot password?</a>
          </div>
          <div class="input-container">
            <div class="input-icon">
              <svg class="input-svg" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path>
              </svg>
            </div>
            <input
              id="password"
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              class="input-field"
              placeholder="••••••••"
              autocomplete="current-password"
              required
              @focus="clearError('password')"
            />
            <button
              type="button"
              class="password-toggle"
              @click="showPassword = !showPassword"
              :aria-label="showPassword ? 'Hide password' : 'Show password'"
            >
              <svg v-if="showPassword" class="toggle-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L6.59 6.59m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21"></path>
              </svg>
              <svg v-else class="toggle-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
              </svg>
            </button>
          </div>
          <p v-if="errors.password" class="error-message">{{ errors.password }}</p>
        </div>

        <!-- Remember me checkbox -->
        <div class="remember-container">
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="form.remember"
              class="checkbox-input"
            />
            <span class="checkbox-custom"></span>
            <span class="checkbox-text">Remember me</span>
          </label>
        </div>

        <!-- Submit button -->
        <button
          type="submit"
          class="submit-btn"
          :disabled="loading"
          :class="{ 'loading-btn': loading }"
        >
          <span v-if="!loading" class="btn-text">Sign In</span>
          <div v-else class="loading-spinner"></div>
          <div class="btn-icon">
            <svg class="btn-svg" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14 5l7 7m0 0l-7 7m7-7H3"></path>
            </svg>
          </div>
        </button>

        <!-- Divider -->
        <div class="divider">
          <span class="divider-text">or continue with</span>
        </div>

        <!-- Social login options -->
        <div class="social-login">
          <button type="button" class="social-btn" @click="handleSocialLogin('google')">
            <svg class="social-icon" viewBox="0 0 24 24" width="24" height="24" xmlns="http://www.w3.org/2000/svg">
              <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
              <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
              <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
              <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
            </svg>
            <span>Google</span>
          </button>
          <button type="button" class="social-btn" @click="handleSocialLogin('github')">
            <svg class="social-icon" viewBox="0 0 24 24" width="24" height="24" xmlns="http://www.w3.org/2000/svg">
              <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z"/>
            </svg>
            <span>GitHub</span>
          </button>
        </div>

        <!-- Sign up link -->
        <p class="signup-text">
          Don't have an account?
          <a href="#" class="signup-link" @click.prevent="$emit('toggleMode')">
            Sign up
          </a>
        </p>
      </form>

      <!-- Success message -->
      <transition name="slide-up">
        <div v-if="successMessage" class="success-message">
          <svg class="success-icon" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
          <p>{{ successMessage }}</p>
        </div>
      </transition>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'

// Emit events for parent components
const emit = defineEmits(['login', 'toggleMode', 'forgotPassword', 'socialLogin'])

// Reactive form data
const form = reactive({
  email: '',
  password: '',
  remember: false
})

// UI state
const showPassword = ref(false)
const loading = ref(false)
const errors = reactive({})
const successMessage = ref('')
const show = ref(false)
const showCard = ref(false)

// Form validation
const validateForm = () => {
  let isValid = true
  errors.email = ''
  errors.password = ''

  // Email validation
  if (!form.email) {
    errors.email = 'Email is required'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.email)) {
    errors.email = 'Please enter a valid email address'
    isValid = false
  }

  // Password validation
  if (!form.password) {
    errors.password = 'Password is required'
    isValid = false
  } else if (form.password.length < 6) {
    errors.password = 'Password must be at least 6 characters'
    isValid = false
  }

  return isValid
}

// Clear error for a specific field
const clearError = (field) => {
  errors[field] = ''
}

// Handle login submission
const handleLogin = async () => {
  if (!validateForm()) return

  loading.value = true
  successMessage.value = ''

  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // Emit login event with form data
    emit('login', { ...form })
    
    successMessage.value = 'Login successful! Redirecting...'
    
    // In a real app, you would redirect or update state here
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  } catch (error) {
    errors.email = 'Login failed. Please check your credentials.'
  } finally {
    loading.value = false
  }
}

// Handle social login
const handleSocialLogin = (provider) => {
  emit('socialLogin', provider)
}

// Handle forgot password
const handleForgotPassword = () => {
  emit('forgotPassword', form.email)
}

// Animation on mount
onMounted(() => {
  setTimeout(() => {
    show.value = true
  }, 100)
  
  setTimeout(() => {
    showCard.value = true
  }, 300)
})
</script>

<style scoped>
/* Container styling */
.login-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #f5f7fa 0%, #e4f0e8 100%);
  padding: 1rem;
  opacity: 0;
  transition: opacity 0.8s ease-out;
  position: relative;
  overflow: hidden;
}

.login-hidden {
  opacity: 0;
}

.login-container {
  opacity: 1;
}

/* Background shapes animation */
.bg-shapes {
  position: absolute;
  width: 100%;
  height: 100%;
  z-index: 0;
}

.shape {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  animation: float 20s infinite linear;
}

.shape-1 {
  width: 300px;
  height: 300px;
  top: -150px;
  left: -100px;
  animation-delay: 0s;
}

.shape-2 {
  width: 200px;
  height: 200px;
  top: 20%;
  right: -80px;
  animation-delay: 5s;
  background: linear-gradient(135deg, #34d399 0%, #10b981 100%);
}

.shape-3 {
  width: 150px;
  height: 150px;
  bottom: 10%;
  left: 5%;
  animation-delay: 10s;
  background: linear-gradient(135deg, #059669 0%, #047857 100%);
}

.shape-4 {
  width: 100px;
  height: 100px;
  top: 40%;
  left: 20%;
  animation-delay: 15s;
  background: linear-gradient(135deg, #10b981 0%, #34d399 100%);
}

.shape-5 {
  width: 250px;
  height: 250px;
  bottom: -100px;
  right: -50px;
  animation-delay: 20s;
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg);
  }
  33% {
    transform: translateY(-20px) rotate(120deg);
  }
  66% {
    transform: translateY(20px) rotate(240deg);
  }
}

/* Login card styling */
.login-card {
  background: white;
  border-radius: 1.5rem;
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.1), 0 0 0 1px rgba(0, 0, 0, 0.05);
  padding: 2.5rem;
  width: 100%;
  max-width: 420px;
  z-index: 10;
  position: relative;
  transform: translateY(20px);
  opacity: 0;
  transition: all 0.6s cubic-bezier(0.34, 1.56, 0.64, 1);
}

.card-hidden {
  transform: translateY(20px);
  opacity: 0;
}

.login-card {
  transform: translateY(0);
  opacity: 1;
}

/* Header styling */
.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.logo-container {
  margin-bottom: 1.5rem;
}

.logo-circle {
  width: 70px;
  height: 70px;
  margin: 0 auto 1rem;
  border-radius: 50%;
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 20px rgba(16, 185, 129, 0.3);
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    box-shadow: 0 10px 20px rgba(16, 185, 129, 0.3);
  }
  50% {
    box-shadow: 0 10px 30px rgba(16, 185, 129, 0.5);
  }
}

.logo-icon {
  width: 32px;
  height: 32px;
  color: white;
}

.login-title {
  font-size: 1.875rem;
  font-weight: 700;
  color: #111827;
  margin-bottom: 0.5rem;
}

.login-subtitle {
  color: #6b7280;
  font-size: 0.95rem;
}

/* Form styling */
.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.input-group {
  transition: all 0.3s ease;
}

.input-label {
  display: block;
  font-size: 0.875rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.5rem;
}

.input-container {
  position: relative;
  display: flex;
  align-items: center;
  border: 1px solid #d1d5db;
  border-radius: 0.75rem;
  transition: all 0.3s ease;
  overflow: hidden;
  background: #f9fafb;
}

.input-container:focus-within {
  border-color: #10b981;
  box-shadow: 0 0 0 3px rgba(16, 185, 129, 0.1);
  background: white;
  transform: translateY(-2px);
}

.input-error .input-container {
  border-color: #ef4444;
}

.input-icon {
  padding: 0 0.875rem;
  display: flex;
  align-items: center;
}

.input-svg {
  width: 1.25rem;
  height: 1.25rem;
  color: #9ca3af;
}

.input-container:focus-within .input-svg {
  color: #10b981;
}

.input-field {
  flex: 1;
  padding: 0.875rem 0;
  padding-right: 3rem;
  border: none;
  background: transparent;
  font-size: 1rem;
  color: #111827;
  outline: none;
}

.input-field::placeholder {
  color: #9ca3af;
}

.password-toggle {
  position: absolute;
  right: 0.875rem;
  background: none;
  border: none;
  cursor: pointer;
  padding: 0.25rem;
  color: #9ca3af;
  transition: color 0.2s;
}

.password-toggle:hover {
  color: #10b981;
}

.toggle-icon {
  width: 1.25rem;
  height: 1.25rem;
}

.error-message {
  color: #ef4444;
  font-size: 0.875rem;
  margin-top: 0.5rem;
  animation: shake 0.5s ease;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}

/* Remember me checkbox */
.remember-container {
  margin-top: -0.5rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 0.875rem;
  color: #6b7280;
  user-select: none;
}

.checkbox-input {
  display: none;
}

.checkbox-custom {
  width: 1.25rem;
  height: 1.25rem;
  border: 2px solid #d1d5db;
  border-radius: 0.375rem;
  margin-right: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s;
}

.checkbox-input:checked + .checkbox-custom {
  background-color: #10b981;
  border-color: #10b981;
}

.checkbox-input:checked + .checkbox-custom::after {
  content: '';
  width: 0.5rem;
  height: 0.5rem;
  background-color: white;
  border-radius: 1px;
}

.checkbox-text {
  font-size: 0.875rem;
}

.forgot-link {
  font-size: 0.875rem;
  color: #10b981;
  text-decoration: none;
  font-weight: 600;
  transition: color 0.2s;
}

.forgot-link:hover {
  color: #059669;
  text-decoration: underline;
}

/* Submit button */
.submit-btn {
  background: linear-gradient(135deg, #10b981 0%, #059669 100%);
  color: white;
  border: none;
  border-radius: 0.75rem;
  padding: 1rem;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  overflow: hidden;
  margin-top: 0.5rem;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(16, 185, 129, 0.3);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.loading-btn {
  cursor: wait;
}

.btn-text {
  margin-right: 0.5rem;
}

.btn-icon {
  width: 1.25rem;
  height: 1.25rem;
  transition: transform 0.3s ease;
}

.submit-btn:hover:not(:disabled) .btn-icon {
  transform: translateX(5px);
}

.btn-svg {
  width: 100%;
  height: 100%;
}

.loading-spinner {
  width: 1.5rem;
  height: 1.5rem;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* Divider */
.divider {
  display: flex;
  align-items: center;
  text-align: center;
  margin: 1.5rem 0;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  border-bottom: 1px solid #e5e7eb;
}

.divider-text {
  padding: 0 1rem;
  color: #6b7280;
  font-size: 0.875rem;
}

/* Social login */
.social-login {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
  margin-bottom: 1.5rem;
}

.social-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.75rem;
  padding: 0.75rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.75rem;
  background: white;
  color: #374151;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
}

.social-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  border-color: #10b981;
}

.social-icon {
  width: 1.25rem;
  height: 1.25rem;
}

/* Sign up link */
.signup-text {
  text-align: center;
  color: #6b7280;
  font-size: 0.95rem;
  margin-top: 1rem;
}

.signup-link {
  color: #10b981;
  font-weight: 600;
  text-decoration: none;
  margin-left: 0.5rem;
  transition: color 0.2s;
}

.signup-link:hover {
  color: #059669;
  text-decoration: underline;
}

/* Success message */
.success-message {
  position: fixed;
  bottom: 2rem;
  left: 50%;
  transform: translateX(-50%);
  background: #10b981;
  color: white;
  padding: 1rem 1.5rem;
  border-radius: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  box-shadow: 0 10px 25px rgba(16, 185, 129, 0.3);
  animation: slideIn 0.5s ease;
  z-index: 100;
}

.success-icon {
  width: 1.25rem;
  height: 1.25rem;
}

@keyframes slideIn {
  from {
    transform: translateX(-50%) translateY(100%);
    opacity: 0;
  }
  to {
    transform: translateX(-50%) translateY(0);
    opacity: 1;
  }
}

/* Slide up transition */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: all 0.5s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  opacity: 0;
  transform: translateY(20px);
}

/* Responsive adjustments */
@media (max-width: 480px) {
  .login-card {
    padding: 2rem 1.5rem;
  }
  
  .social-login {
    grid-template-columns: 1fr;
  }
}
</style>