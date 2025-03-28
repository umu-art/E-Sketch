import React from 'react';
import { fireEvent, render, screen } from '@testing-library/react';
import BoardCard from './BoardCard';
import { useNavigate } from 'react-router-dom';

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

    const MockMeta = ({ title, description, avatar }) => (
        <div id="Meta" className="ant-card-meta">
            {avatar}
            <div className="ant-card-meta-title">{title}</div>
            <div className="ant-card-meta-description">{description}</div>
        </div>
    );

    const MockCard = ({ className, style, actions, cover, children }) => (
        <div id="Card" className={className} style={style}>
            {cover}
            <div id="children">{children}</div>
            <div id="actions">{actions}</div>
        </div>
    );
    MockCard.Meta = MockMeta;

    return {
        ...antd,
        Avatar: ({ src }) => <img id="Avatar" src={src} alt="avatar" />,
        Card: MockCard,
        Image: ({ src, onClick }) => <img id="Image" src={src} onClick={onClick} alt="cover" />,
        Modal: ({ children, open, onCancel, title }) => (
            open ? (
                <div id="Modal">
                    <h3>{title}</h3>
                    <button onClick={onCancel}>Close Modal</button>
                    {children}
                </div>
            ) : null
        ),
    };
});

jest.mock('@ant-design/icons', () => ({
    EditOutlined: ({ onClick }) => <button onClick={onClick}>Edit</button>,
    EllipsisOutlined: () => <button>More</button>,
    SettingOutlined: () => <button>Settings</button>,
}));

jest.mock('../UpdateBoardForm/UpdateBoardForm', () => 
    ({ data, onDataChange, closeModal }) => (
        <div>
            <p>Update Board Form</p>
            <p>Current Name: {data.name}</p>
            <button onClick={() => onDataChange({ ...data, name: 'Updated Board' })}>
                Update Data
            </button>
            <button onClick={closeModal}>Submit Form</button>
        </div>
    )
);

describe('BoardCard Component', () => {
    const mockNavigate = jest.fn();
    const mockBoardData = {
        id: '123',
        name: 'Test Board',
        description: 'Test Description',
        preview: 'preview-token',
        editable: true,
        ownerInfo: {
            id: 'owner-123',
        },
    };

    beforeEach(() => {
        useNavigate.mockReturnValue(mockNavigate);
    });

    afterEach(() => {
        jest.clearAllMocks();
    });

    it('renders empty state when no data is provided', () => {
        render(
            <BoardCard />
        );

        expect(screen.getByAltText('cover')).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Settings' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'Edit' })).toBeInTheDocument();
        expect(screen.getByRole('button', { name: 'More' })).toBeInTheDocument();
    });

    it('renders with board data', () => {
        render(
            <BoardCard data={mockBoardData} />
        );

        expect(screen.getByAltText('cover')).toHaveAttribute(
        'src',
        `/preview?boardId=123&token=preview-token`
        );
        expect(screen.getByText('Test Board')).toBeInTheDocument();
        expect(screen.getByText('Test Description')).toBeInTheDocument();
        expect(screen.getByAltText('avatar')).toHaveAttribute(
        'src',
        'https://api.dicebear.com/7.x/miniavs/svg?seed=owner-123'
        );
    });

    it('navigates to board when cover image is clicked', () => {
        render(
            <BoardCard data={mockBoardData} />
        );

        fireEvent.click(screen.getByAltText('cover'));
        expect(mockNavigate).toHaveBeenCalledWith('/app/board/123');
    });

    it('opens update modal when edit button is clicked', () => {
        render(
            <BoardCard data={mockBoardData} />
        );

        fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
        expect(screen.getByText('Изменить доску')).toBeInTheDocument();
        expect(screen.getByText('Update Board Form')).toBeInTheDocument();
    });

    it('closes update modal when close button is clicked', () => {
        render(
            <BoardCard data={mockBoardData} />
        );

        fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
        expect(screen.getByText('Изменить доску')).toBeInTheDocument();
        
        fireEvent.click(screen.getByRole('button', { name: 'Close Modal' }));
        expect(screen.queryByText('Изменить доску')).not.toBeInTheDocument();
    });

    it('updates board data through the update form', () => {
        render(
            <BoardCard data={mockBoardData} />
        );

        fireEvent.click(screen.getByRole('button', { name: 'Edit' }));
        expect(screen.getByText('Current Name: Test Board')).toBeInTheDocument();
        
        fireEvent.click(screen.getByRole('button', { name: 'Update Data' }));
        expect(screen.getByText('Current Name: Updated Board')).toBeInTheDocument();
    });

    it('hides action buttons when editable is false', () => {
        render(
            <BoardCard data={mockBoardData} editable={false}/>
        );

        expect(screen.queryByRole('button', { name: 'Settings' })).not.toBeInTheDocument();
        expect(screen.queryByRole('button', { name: 'Edit' })).not.toBeInTheDocument();
        expect(screen.queryByRole('button', { name: 'More' })).not.toBeInTheDocument();
    });
});