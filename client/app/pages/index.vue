<template>
  <div class="app-container">
    <LoginForm 
      @login="handleLogin"
      @toggleMode="toggleMode"
      @forgotPassword="handleForgotPassword"
      @socialLogin="handleSocialLogin"
    />
  </div>
</template>

<script setup>
// This is a demo page showing how to use the LoginForm component

const handleLogin = (formData) => {
  console.log('Login attempted with:', formData)
  alert(`Login submitted with email: ${formData.email}`)
  // In a real app, you would make an API call here
}

const toggleMode = () => {
  console.log('Switching to signup mode')
  alert('Switching to signup form (you would implement this)')
}

const handleForgotPassword = (email) => {
  console.log('Forgot password for:', email)
  alert(`Password reset requested for: ${email}`)
}

const handleSocialLogin = (provider) => {
  console.log('Social login with:', provider)
  alert(`Social login with ${provider} (you would implement OAuth)`)
}
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

.app-container {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}
</style>