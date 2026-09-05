// =============================================================================
// MovieHub Production CI/CD Pipeline (Jenkinsfile)
// Automates: Linting -> Testing -> Building -> Deploying to EC2 -> Live Smoke Test
// =============================================================================

pipeline {
    agent any

    options {
        timestamps()
        timeout(time: 20, unit: 'MINUTES')
        disableConcurrentBuilds()
        ansiColor('xterm')
    }

    environment {
        // App and Environment Metadata
        APP_NAME         = 'moviehub'
        DOMAIN           = 'https://moviehub.nostackdev.online'
        DOCKER_REGISTRY  = 'docker.io'
        DOCKER_ORG       = 'yourdockerhubuser' // Change to your Docker Hub username if pushing images
        DOCKER_CREDS_ID  = 'dockerhub-credentials'
        
        // Dynamic build identifiers
        GIT_COMMIT_SHORT = ''
        IMAGE_TAG        = ''
    }

    stages {
        // ── 1. SCM Checkout & Versioning ──────────────────────────────────────
        stage('1. Checkout & Versioning') {
            steps {
                script {
                    GIT_COMMIT_SHORT = sh(script: 'git rev-parse --short HEAD', returnStdout: true).trim()
                    IMAGE_TAG = "v${BUILD_NUMBER}-${GIT_COMMIT_SHORT}"
                    currentBuild.displayName = "#${BUILD_NUMBER} (${GIT_COMMIT_SHORT})"

                    echo """
                    =============================================================
                    🚀 Triggered Build: #${BUILD_NUMBER}
                    🏷️ Commit:          ${GIT_COMMIT_SHORT}
                    🌐 Live Domain:     ${DOMAIN}
                    =============================================================
                    """
                }
            }
        }

        // ── 2. Run Backend & Frontend Tests ───────────────────────────────────
        stage('2. Test & Quality Checks') {
            parallel {
                stage('Go Backend Tests') {
                    steps {
                        dir('moviehub') {
                            echo "🧪 Running Go static analysis and unit tests..."
                            sh 'go vet ./...'
                            sh 'go test -v -race ./internal/...'
                        }
                    }
                }
                stage('Frontend Build Test') {
                    steps {
                        dir('frontend') {
                            echo "🧪 Validating frontend dependencies & production build..."
                            sh 'npm ci'
                            sh 'npm run build'
                        }
                    }
                }
            }
        }

        // ── 3. Build Docker Images ────────────────────────────────────────────
        stage('3. Build Containers') {
            steps {
                script {
                    echo "🐳 Building microservice & frontend Docker images..."
                    dir('moviehub') {
                        // Builds all containers defined in docker-compose.prod.yml
                        sh "docker compose -f docker-compose.prod.yml build"
                    }
                }
            }
        }

        // ── 4. Zero-Downtime Deployment ───────────────────────────────────────
        stage('4. Deploy to Production') {
            steps {
                script {
                    echo "🚢 Applying rolling update with Docker Compose..."
                    dir('moviehub') {
                        // Replaces running containers with new builds without taking down databases
                        sh "docker compose -f docker-compose.prod.yml up -d --remove-orphans"
                    }
                }
            }
        }

        // ── 5. Live Smoke Test & Health Verification ──────────────────────────
        stage('5. Smoke Test & Health Check') {
            steps {
                script {
                    echo "🩺 Verifying production health at ${DOMAIN}..."
                    dir('moviehub') {
                        sh """
                            # Run automated health and sanity checks against the live domain
                            python3 scripts/test_and_seed.py --base-url ${DOMAIN}
                        """
                    }
                    echo "✅ Live deployment verified and operational!"
                }
            }
        }
    }

    // ── Post Actions & Cleanup ────────────────────────────────────────────────
    post {
        always {
            echo "🧹 Pruning dangling Docker builder cache..."
            sh "docker image prune -f || true"
        }
        success {
            echo """
            =============================================================
            🎉 SUCCESS: Commit ${GIT_COMMIT_SHORT} is LIVE on production!
            🌍 Live URL: ${DOMAIN}
            =============================================================
            """
        }
        failure {
            echo """
            =============================================================
            ❌ BUILD FAILED: Deployment or tests did not succeed.
            Check the console output above for error logs.
            =============================================================
            """
        }
    }
}
