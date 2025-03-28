import React from 'react';
import { render, screen } from '@testing-library/react';
import { BrowserRouter as Router } from 'react-router-dom';
import SignUpPage from './SignUpPage';

jest.mock('./SignUpForm', () => () => <div data-testid="signup-form" />);
jest.mock('../../FormPage/FormPage', () => ({ children, form, title }) => (
    <div>
        <h1 data-testid="form-page-title">{title}</h1>
        <div data-testid="form-page-children">{children}</div>
        <div>{form}</div>
    </div>
));

describe('SignUpPage', () => {
    it('renders correctly with all components', () => {
        render(
            <Router>
                <SignUpPage />
            </Router>
        );

        expect(screen.getByTestId('form-page-title')).toHaveTextContent('Регистрация');
        
        expect(screen.getByTestId('signup-form')).toBeInTheDocument();
        
        expect(screen.getByText('Уже есть аккаунт?')).toBeInTheDocument();
        expect(screen.getByRole('link', { name: 'Вход' })).toBeInTheDocument();
        expect(screen.getByRole('link', { name: 'Вход' })).toHaveAttribute('href', 'signin');
    });
});