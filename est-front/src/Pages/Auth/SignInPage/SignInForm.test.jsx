import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import SignInForm from './SignInForm';
import { BrowserRouter as Router, useNavigate } from 'react-router-dom';
import { message } from 'antd';

window.matchMedia = window.matchMedia || function() {
    return {
        matches: false,
        addListener: function() {},
        removeListener: function() {}
    };
};

jest.mock('react-router-dom', () => ({
    ...jest.requireActual('react-router-dom'),
    useNavigate: jest.fn(),
}));

jest.mock('antd', () => {
    const antd = jest.requireActual('antd');
    return {
        ...antd,
        message: {
            useMessage: jest.fn(),
        },
    };
});

jest.mock('est_proxy_api', () => ({
    UserApi: jest.fn().mockImplementation(() => ({
        apiClient: {
            basePath: '',
            defaultHeaders: {},
        },
        login: jest.fn(),
    })),
}));

describe('SignInForm', () => {
    const redirectTo = '/dashboard';
    let mockNavigate;
    let mockMessageApi;
    let mockLogin;

    beforeEach(() => {
        mockNavigate = jest.fn();
        useNavigate.mockReturnValue(mockNavigate);

        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockLogin = jest.fn();
        const mockUserApi = {
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            login: mockLogin,
        };
        require('est_proxy_api').UserApi.mockImplementation(() => mockUserApi);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });
    
    const renderComponent = (props = {}) => {
        return render(
            <Router>
                <SignInForm redirectTo={redirectTo} {...props} />
            </Router>
        );
    };

    it('renders the form correctly', () => {
        renderComponent();

        expect(screen.getByLabelText('Почта')).toBeInTheDocument();
        expect(screen.getByLabelText('Пароль')).toBeInTheDocument();
        expect(screen.getByRole('checkbox', { name: 'Запомнить меня?' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Войти' })).toBeInTheDocument();
        expect(screen.getByText('Забыли пароль?')).toBeInTheDocument();
    });

    it('should validate required fields before submission', async () => {
        renderComponent();

        fireEvent.click(screen.getByRole('button', { name: 'Войти' }));

        await waitFor(async () => {
            expect(await screen.findAllByText(/Пожалуйста, введите/i)).toHaveLength(2);
            expect(mockLogin).not.toHaveBeenCalled();
        });
    });

    it('should validate email format', async () => {
        renderComponent();
        
        fireEvent.change(screen.getByLabelText('Почта'), { target: { value: 'invalid-email' } });
        fireEvent.click(screen.getByRole('button', { name: 'Войти' }));

        await waitFor(async () => {
            expect(await screen.findByText('Некорректная почта!')).toBeInTheDocument();
            expect(mockLogin).not.toHaveBeenCalled();
        });
    });

    it('should successfully submit form with valid data', async () => {
        mockLogin.mockResolvedValue({});
        renderComponent();

        fireEvent.change(screen.getByLabelText('Почта'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText('Пароль'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: 'Войти' }));

        await waitFor(() => {
            expect(mockLogin).toHaveBeenCalledWith(expect.objectContaining({
                authDto: {
                    email: 'test@example.com',
                    passwordHash: 'password123',
                }
            }));
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'success',
                content: 'Авторизация прошла успешно!',
            });
            expect(mockNavigate).toHaveBeenCalledWith(redirectTo);
        });
    });

    it('should handle login error from API', async () => {
        const errorResponse = {
            response: {
                text: 'Неверные учетные данные',
            },
        };
        mockLogin.mockRejectedValueOnce(errorResponse);
        renderComponent();

        fireEvent.change(screen.getByLabelText('Почта'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText('Пароль'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: 'Войти' }));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Неверные учетные данные',
            });
            expect(mockNavigate).not.toHaveBeenCalled();
        });
    });

    it('should handle unexpected errors', async () => {
        mockLogin.mockRejectedValueOnce(new Error('Network error'));
        renderComponent();

        fireEvent.change(screen.getByLabelText('Почта'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText('Пароль'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: 'Войти' }));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith(expect.objectContaining({
                type: 'error',
                content: 'Произошла ошибка, попробуйте позже',
            }));
        });
    });
});