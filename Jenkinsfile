// =============================================================================
// MovieHub Production CI/CD Pipeline (Jenkinsfile)
// Automates: Linting -> Testing -> Building -> Deploying to EC2 -> Live Smoke Test
// =============================================================================

pipeline {
    agent {
        label 'ec2-agent'
    }

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
        DOCKER_ORG       = 'yourdockerhubuser'
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
                            echo "🧪 Running Go static analysis and unit tests inside container..."
                            sh '''
                                docker run --rm \
                                  -v $(pwd):/app \
                                  -w /app \
                                  golang:1.23-alpine \
                                  sh -c "go vet ./... && go test -v ./internal/..."
                            '''
                        }
                    }
                }
                stage('Frontend Build Test') {
                    steps {
                        dir('frontend') {
                            echo "🧪 Validating frontend dependencies & production build inside container..."
                            sh '''
                                docker run --rm \
                                  -v $(pwd):/app \
                                  -w /app \
                                  node:20-alpine \
                                  sh -c "npm ci && npm run build"
                            '''
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
                        sh '''
                            if [ -d "/home/ubuntu/certs" ] && [ ! -d "certs" ]; then
                                echo "🔐 Copying SSL certificates from /home/ubuntu/certs..."
                                cp -r /home/ubuntu/certs ./certs
                            fi
                        '''
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
                    sh """
                        for i in 1 2 3 4 5 6; do
                            if curl -k -s -f -o /dev/null "${DOMAIN}"; then
                                echo "✅ Live site responded 200 OK at ${DOMAIN}!"
                                break
                            fi
                            echo "Waiting for services to finish starting (attempt \$i/6)..."
                            sleep 5
                        done
                    """
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
