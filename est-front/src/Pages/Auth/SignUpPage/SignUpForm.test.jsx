import React from 'react';
import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import SignUpForm from './SignUpForm';
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
        register: jest.fn(),
    })),
}));

describe('SignUpForm', () => {
    let mockNavigate;
    let mockMessageApi;
    let mockRegister;

    beforeEach(() => {
        mockNavigate = jest.fn();
        useNavigate.mockReturnValue(mockNavigate);

        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockRegister = jest.fn();
        const mockUserApi = {
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            register: mockRegister,
        };
        require('est_proxy_api').UserApi.mockImplementation(() => mockUserApi);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });
    
    const renderComponent = (props = {}) => {
        return render(
            <Router>
                <SignUpForm />
            </Router>
        );
    };

    it('renders the form correctly', () => {
        renderComponent();

        expect(screen.getByLabelText('Имя пользователя')).toBeInTheDocument();
        expect(screen.getByLabelText('Email')).toBeInTheDocument();
        expect(screen.getByLabelText('Пароль')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Зарегистрироваться' })).toBeInTheDocument();
    });

    it('should validate required fields before submission', async () => {
        renderComponent();

        fireEvent.click(screen.getByRole('button', { name: 'Зарегистрироваться' }));

        await waitFor(async () => {
            expect(await screen.findAllByText(/Пожалуйста, введите/i)).toHaveLength(3);
            expect(mockRegister).not.toHaveBeenCalled();
        });
    });

    it('should validate email format', async () => {
        renderComponent();
        
        fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'invalid-email' } });
        fireEvent.click(screen.getByRole('button', { name: 'Зарегистрироваться' }));

        await waitFor(async () => {
            expect(await screen.findByText('Некорректный email!')).toBeInTheDocument();
            expect(mockRegister).not.toHaveBeenCalled();
        })
    });

    it('should successfully submit form with valid data', async () => {
        mockRegister.mockResolvedValue({});
        renderComponent();

        fireEvent.change(screen.getByLabelText('Имя пользователя'), { target: { value: 'testuser' } });
        fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText('Пароль'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button'));

        await waitFor(() => {
            expect(mockRegister).toHaveBeenCalledWith(expect.objectContaining({
                registerDto: {
                    email: 'test@example.com',
                    passwordHash: 'password123',
                    username: 'testuser',
                }
            }));
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'success',
                content: 'Регистрация прошла успешно!',
            });
            expect(mockNavigate).toHaveBeenCalledWith('/auth/confirm');
        });
    });

    it('should handle registration error from API', async () => {
        const errorResponse = {
            response: {
                text: 'Email уже используется',
            },
        };
        mockRegister.mockRejectedValueOnce(errorResponse);
        renderComponent();

        fireEvent.change(screen.getByLabelText('Имя пользователя'), { target: { value: 'testuser' } });
        fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText('Пароль'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button', { name: 'Зарегистрироваться' }));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Email уже используется',
            });
            expect(mockNavigate).not.toHaveBeenCalled();
        });
    });

    it('should handle unexpected errors', async () => {
        mockRegister.mockRejectedValueOnce(new Error('Network error'));
        renderComponent();

        fireEvent.change(screen.getByLabelText('Имя пользователя'), { target: { value: 'testuser' } });
        fireEvent.change(screen.getByLabelText('Email'), { target: { value: 'test@example.com' } });
        fireEvent.change(screen.getByLabelText('Пароль'), { target: { value: 'password123' } });

        fireEvent.click(screen.getByRole('button'));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith(expect.objectContaining({
                type: 'error',
                content: 'Произошла ошибка, попробуйте позже',
            }));
        });
    });
});