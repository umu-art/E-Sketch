import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import { BrowserRouter as Router, useLocation, useNavigate } from 'react-router-dom';
import EmailConfirmPage from './EmailConfirmPage';
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
    useLocation: jest.fn(),
    useNavigate: jest.fn(),
}));

jest.mock('antd', () => {
    const antd = jest.requireActual('antd');
    return {
        ...antd,
        message: {
            useMessage: jest.fn(),
        },
        Link: ({ to, children }) => <a href={to}>{children}</a>,
    };
});

jest.mock('est_proxy_api', () => ({
    UserApi: jest.fn().mockImplementation(() => ({
        getSelf: jest.fn(),
        logout: jest.fn(),
    })),
}));

describe('EmailConfirmPage', () => {
    let mockNavigate;
    let mockMessageApi;
    let mockConfirm;
    let mockLocation;

    beforeEach(() => {
        mockNavigate = jest.fn();
        useNavigate.mockReturnValue(mockNavigate);

        mockMessageApi = {
            open: jest.fn(),
        };
        message.useMessage.mockReturnValue([mockMessageApi, jest.fn()]);

        mockConfirm = jest.fn();
        const mockUserApi = {
            apiClient: {
                basePath: '',
                defaultHeaders: {},
            },
            confirm: mockConfirm,
        };
        require('est_proxy_api').UserApi.mockImplementation(() => mockUserApi);

        mockLocation = {
            search: '',
        };
        useLocation.mockReturnValue(mockLocation);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    const renderComponent = (search = '') => {
        mockLocation.search = search;
        return render(
            <Router>
                <EmailConfirmPage />
            </Router>
        );
    };

    it('renders correctly without token', () => {
        renderComponent();
        
        expect(screen.getByText('Подтверждение почты')).toBeInTheDocument();
        expect(screen.getByText('Письмо с подтверждением отправлено на вашу электронную почту.')).toBeInTheDocument();
        expect(screen.getByLabelText('Код')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Подтвердить' })).toBeInTheDocument();
        expect(screen.getByText(/Если вы не получили письма с подтверждением/i)).toBeInTheDocument();
    });

    it('calls confirm API when token is present in URL', async () => {
        mockConfirm.mockResolvedValue({});
        renderComponent('?token=test-token');

        await waitFor(() => {
            expect(mockConfirm).toHaveBeenCalledWith({
                confirmationDto: {
                    token: 'test-token',
                }
            });
        });
    });

    it('handles successful confirmation', async () => {
        mockConfirm.mockResolvedValue({});
        renderComponent();
        
        fireEvent.change(screen.getByLabelText('Код'), { target: { value: 'test-token' } });
        fireEvent.click(screen.getByRole('button', { name: 'Подтвердить' }));

        await waitFor(() => {
            expect(mockConfirm).toHaveBeenCalledWith({
                confirmationDto: {
                    token: 'test-token',
                }
            });
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'success',
                content: 'Регистрация прошла успешно!'
            });
            expect(mockNavigate).toHaveBeenCalledWith('/app');
        });
    });

    it('handles confirmation error with response', async () => {
        const errorResponse = {
            response: {
                text: 'Invalid token',
            },
        };
        mockConfirm.mockRejectedValue(errorResponse);
        renderComponent();
        
        fireEvent.change(screen.getByLabelText('Код'), { target: { value: 'invalid-token' } });
        fireEvent.click(screen.getByRole('button', { name: 'Подтвердить' }));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Invalid token',
            });
        });
    });

    it('handles unexpected errors', async () => {
        mockConfirm.mockRejectedValue(new Error('Network error'));
        renderComponent();
        
        fireEvent.change(screen.getByLabelText('Код'), { target: { value: 'test-token' } });
        fireEvent.click(screen.getByRole('button', { name: 'Подтвердить' }));

        await waitFor(() => {
            expect(mockMessageApi.open).toHaveBeenCalledWith({
                type: 'error',
                content: 'Что-то пошло не так ;(',
            });
        });
    });

    it('updates token state when input changes', () => {
        renderComponent();
        
        const input = screen.getByLabelText('Код');
        fireEvent.change(input, { target: { value: 'new-token' } });
        
        expect(input).toHaveValue('new-token');
    });
});